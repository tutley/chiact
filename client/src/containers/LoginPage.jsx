import React, { PropTypes } from 'react';
import Auth from '../modules/Auth';
import LoginForm from '../components/LoginForm.jsx';
import Axios from 'axios';

class LoginPage extends React.Component {

  /**
   * Class constructor.
   */
  constructor(props, context) {
    super(props, context);

    const storedMessage = localStorage.getItem('successMessage');
    let successMessage = '';

    if (storedMessage) {
      successMessage = storedMessage;
      localStorage.removeItem('successMessage');
    }

    // set the initial component state
    this.state = {
      errors: {},
      successMessage,
      user: {
        email: '',
        password: ''
      }
    };

    this.processForm = this.processForm.bind(this);
    this.changeUser = this.changeUser.bind(this);
  }

  /**
   * Process the form.
   *
   * @param {object} event - the JavaScript event object
   */
  processForm(event) {
    // prevent default action. in this case, action is the form submission event
    event.preventDefault();

    var authOptions = {
        method: 'GET',
        url: '/api/1/login',
        headers: {
            'Authorization': 'Basic ' + (new Buffer(this.state.user.email + ':' + this.state.user.password).toString('base64')),
            'Content-Type': 'application/x-www-form-urlencoded'
        },
        json: true
    };

    Axios(authOptions).then(response => {
      this.setState({
        errors: {}
      });
      Auth.authenticateUser(response.data.token);
      console.log(response.data.token);
      this.context.router.replace('/');
    })
    .catch(error => {
      console.log(error);
      const errors = error ? error : {};
      errors.summary = error.message;
      this.setState({
        errors
      });
    });
  }

  /**
   * Change the user object.
   *
   * @param {object} event - the JavaScript event object
   */
  changeUser(event) {
    const field = event.target.name;
    const user = this.state.user;
    user[field] = event.target.value;

    this.setState({
      user
    });
  }

  /**
   * Render the component.
   */
  render() {
    return (
      <LoginForm
        onSubmit={this.processForm}
        onChange={this.changeUser}
        errors={this.state.errors}
        successMessage={this.state.successMessage}
        user={this.state.user}
      />
    );
  }

}

LoginPage.contextTypes = {
  router: PropTypes.object.isRequired
};

export default LoginPage;
