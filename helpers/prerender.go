package helpers


import (
	"net/http"

	"log"
	"github.com/tutley/chiact/models"
	"regexp"
	"github.com/tutley/phantomgo"
	"io/ioutil"
	"time"
	"gopkg.in/mgo.v2"
)


func renderAndSave(path string, sess *mgo.Session, dbName string) {
	p := &phantomgo.Param{
		Method:       "GET",
		Url:          "http://127.0.0.1:3333"+path,
		Header:       http.Header{},
		UsePhantomJS: true,
	}

	browser := phantomgo.NewPhantom()
	resp, err := browser.Download(p)
	if err != nil {
		log.Println("There was a problem with PhantomJS - ", err)
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println("Error reading the data returned from PhantomJS - ", err)
	}
	bodyString := string(body)
	// Now lets put it in the DB
	np := models.Page{
		URL: path,
		Modified: time.Now(),
		Content: bodyString,
	}
	db := sess.DB(dbName)
	err = np.Save(db)
	if err != nil {
		log.Println("Error, page not saved to database - ", err)
	}
	sess.Close()
}

// PrerenderMiddleware checks the database to see if the page has already been
// prerendered, and if so it serves that string as the response. If not, it
// forwards the request on to the main app, and also launches a goroutine to
// get a parse of the page and put it into the main db
func PrerenderMiddleware (next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		path := r.RequestURI
		re := regexp.MustCompile(".(gif|jpg|jpeg|tiff|png|js|ico)|(css/)|(js/)")
		found := re.FindString(path)
		if len(found) < 1 {
			db := GetDb(r.Context())
			if db == nil {
				log.Print("No database context")
				http.Error(w, "Not authorized", 401)
			}
			var pageExists bool = true
			var pageIsGood bool = false
			page, err := models.FindPageByURL(path, db)
			if err != nil {
				log.Print("Error finding prerendered page ", err)
				pageExists = false
			} else {
				pageIsGood = true
				if page == nil {
					pageExists = false
				}
				if len(page.Content) < 1 {
					pageExists = false
				}
				// if the page is more than a certain time period old, redo it
				tStart := time.Now().Add(-(time.Minute*60))

				if page.Modified.Before(tStart) {
					log.Println("The page was prerendered more than 60 minutes ago")
					pageIsGood = false
				}
				if pageExists {
					//log.Println("*******************************************************")
					//log.Println("**********     HOLYCRAPITWORKED     *******************")
					//log.Println("*******************************************************")
					log.Println("**** Serving Stored Version from Database ****")
					w.WriteHeader(http.StatusOK)
					w.Write([]byte(page.Content))
				}
			}

			// So here we launch the goroutine to rerender the page
			if r.Host != "127.0.0.1:3333" {
				if !pageExists || !pageIsGood {
					log.Println("@#@ starting phantom prerendering process #@#")
					// we're kicking off a goroutine and moving on, so we have to send
					// along a copy of the database session or it will close prematurely
					go renderAndSave(path, db.Session.Copy(), db.Name)
				}
			}

		}
		next.ServeHTTP(w, r.WithContext(r.Context()))
	});
}