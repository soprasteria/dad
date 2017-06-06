// React
import React from 'react';
import PropTypes from 'prop-types';
import { connect } from 'react-redux';
import { push } from 'react-router-redux';
import DocumentTitle from 'react-document-title';

import AuthPage from '../auth/login.component';

// HomeComponent displaying either the register/login component or information about Dad when authenticated
class HomeComponent extends React.Component {

  redirectIfAuthenticated = (props) => {
    const { isAuthenticated, redirect } = props;
    if (isAuthenticated) {
      redirect('/projects');
    }
  }

  componentWillMount = () => this.redirectIfAuthenticated(this.props);

  componentWillReceiveProps = (nextProps) => this.redirectIfAuthenticated(nextProps);

  render = () => {
    const { isAuthenticated } = this.props;
    if (isAuthenticated) {
      return <div />;
    } else {
      return (
        <DocumentTitle title='D.A.D - Login'>
          <AuthPage />
        </DocumentTitle>
      );
    }
  }
}

HomeComponent.propTypes = {
  isAuthenticated: PropTypes.bool.isRequired,
  redirect: PropTypes.func.isRequired
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
    redirect: (path) => {
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
