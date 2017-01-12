import React from 'react';
import Auth from '../modules/Auth';
import Dashboard from '../components/Dashboard.jsx';
import Axios from 'axios';

class DashboardPage extends React.Component {

  /**
   * Class constructor.
   */
  constructor(props) {
    super(props);

    this.state = {
      name: ''
    };
  }

  /**
   * This method will be executed after initial rendering.
   */
  componentDidMount() {
    Axios({
      method: 'GET',
      url: '/api/1/me',
      headers: {
          'Content-type' : 'application/x-www-form-urlencoded',
          'Authorization': Auth.getToken()
      },
      json: true
    }).then(response => {
      this.setState({
        name: response.data.name
      });
    }).catch(error => {
      console.log(error);
      // do something else?
    });
  }

  /**
   * Render the component.
   */
  render() {
    return (<Dashboard name={this.state.name} />);
  }

}

export default DashboardPage;
