import React, { PropTypes } from 'react';
import SignUpForm from '../components/SignUpForm.jsx';
import Axios from 'axios';

class SignUpPage extends React.Component {

  /**
   * Class constructor.
   */
  constructor(props, context) {
    super(props, context);

    // set the initial component state
    this.state = {
      errors: {},
      user: {
        email: '',
        name: '',
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

    // Use Axios to submit form data
    Axios.post('/api/1/signup', {
      name: this.state.user.name,
      email: this.state.user.email,
      password: this.state.user.password
    })
    .then(response => {
      this.setState({
        errors: {}
      });
      localStorage.setItem('successMessage', response.statusText);
      this.context.router.replace('/login');
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
      <SignUpForm
        onSubmit={this.processForm}
        onChange={this.changeUser}
        errors={this.state.errors}
        user={this.state.user}
      />
    );
  }

}

SignUpPage.contextTypes = {
  router: PropTypes.object.isRequired
};

export default SignUpPage;
