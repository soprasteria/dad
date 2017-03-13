// React
import React from 'react';
import { connect } from 'react-redux';
import { push } from 'react-router-redux';

import AuthPage from '../auth/login.component';

// HomeComponent displaying either the register/login component or information about Dad when authenticated
class HomeComponent extends React.Component {

  redirectIfAuthenticated = (props) => {
    const { isAuthenticated, redirect } = props;
    if (isAuthenticated) {
      redirect('/projects');
    } else {
      document.title = 'D.A.D - Login';
    }
  }

  componentWillMount = () => this.redirectIfAuthenticated(this.props);

  componentWillReceiveProps = (nextProps) => this.redirectIfAuthenticated(nextProps);

  render = () => {
    const { isAuthenticated } = this.props;
    if (isAuthenticated) {
      return <div/>;
    } else {
      return <AuthPage/>;
    }
  }
}

HomeComponent.propTypes = {
  isAuthenticated : React.PropTypes.bool.isRequired,
  redirect: React.PropTypes.func.isRequired
};

// Function to map state to container props
const mapStateToProps = (state) => {
  const { auth } = state;
  const { isAuthenticated } = auth;

  return {
    isAuthenticated
  };
};

// Function to map dispatch to container props
const mapDispatchToProps = (dispatch) => {
  return {
    redirect: path => {
      dispatch(push(path));
    }
  };
};

// Redux container to Sites component
const Home = connect(
  mapStateToProps,
  mapDispatchToProps
)(HomeComponent);

export default Home;
