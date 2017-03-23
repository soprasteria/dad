// React
import React from 'react';
import { connect } from 'react-redux';
import { push } from 'react-router-redux';
import { Button, Container, Form, Header, Message, Segment } from 'semantic-ui-react';
import Joi from 'joi-browser';

import { parseError } from '../../modules/utils/forms';

import AuthThunks from '../../modules/auth/auth.thunk';

import './login.component.scss';

// Signin Pane containing fields to log in the application
class LoginComponent extends React.Component {

  state = { errors: { details: [], fields: {} }, auth: {} }

  schema = Joi.object().keys({
    username: Joi.string().trim().alphanum().required().label('Username'),
    password: Joi.string().trim().min(6).required().label('Password')
  })

  componentWillMount = () => {
    const errorMessage = this.props.errorMessage;
    if (this.props.isAuthenticated && !errorMessage) {
      this.props.redirect(this.props.redirectTo);
    }
    if (errorMessage) {
      this.setState({ errors: { details: [errorMessage], fields:{} }, auth: {} });
    }
  }

  componentWillReceiveProps = (nextProps) => {
    const errorMessage = nextProps.errorMessage;
    if (this.props.isAuthenticated && !errorMessage) {
      this.props.redirect(this.props.redirectTo);
    }
    if (errorMessage) {
      this.setState({ errors: { details: [errorMessage], fields:{} }, auth: {} });
    }
  }

  handleSubmit = (e) => {
    e.preventDefault();
    const { error } = Joi.validate(this.state.auth, this.schema, { abortEarly: false });
    if (error) {
      this.setState({ errors: parseError(error) });
    } else {
      this.props.logUser(this.state.auth);
    }
  }

  handleChange = (e, { name, value }) => {
    const { auth, errors } = this.state;
    const state = {
      auth: { ...auth, [name]: value },
      errors: { details: [...errors.details], fields: { ...errors.fields } }
    };
    delete state.errors.fields[name];
    this.setState(state);
  }

  render = () => {
    const { isFetching } = this.props;
    const { fields, details } = this.state.errors;
    return (
      <Container text>
        <Segment className='login-component' padded raised>
          <Header as='h1'>Login</Header>
          <Form error={Boolean(details.length)} onSubmit={this.handleSubmit}>
            <Form.Input required error={fields['username']} label='Username' onChange={this.handleChange}
              type='text' name='username' autoComplete='off' placeholder='LDAP username'
            />
            <Form.Input required error={fields['password']} label='Password' onChange={this.handleChange}
              type='password' name='password' autoComplete='off' placeholder='Password'
            />
            <Message error list={details}/>
            <Button fluid color='green' content='Login' loading={isFetching} />
          </Form>
        </Segment>
      </Container>
    );
  }
};

LoginComponent.propTypes = {
  isAuthenticated: React.PropTypes.bool.isRequired,
  isFetching: React.PropTypes.bool.isRequired,
  logUser: React.PropTypes.func.isRequired,
  redirect: React.PropTypes.func.isRequired,
  errorMessage: React.PropTypes.string,
  redirectTo : React.PropTypes.string
};

// Function to map state to container props
const mapStateToProps = (state) => {
  const { auth } = state;
  const { isAuthenticated, errorMessage, isFetching } = auth;
  const redirectTo = state.routing.locationBeforeTransitions.query.next || '/';
  return {
    isAuthenticated,
    errorMessage,
    redirectTo,
    isFetching
  };
};

// Function to map dispatch to container props
const mapDispatchToProps = (dispatch) => {
  return {
    logUser: (username, password) => {
      dispatch(AuthThunks.loginUser(username, password));
    },
    redirect: (path) => {
      dispatch(push(path));
    }
  };
};

// Redux container to Sites component
const AuthPage = connect(
  mapStateToProps,
  mapDispatchToProps
)(LoginComponent);

export default AuthPage;
