// React
import React from 'react';
import { IndexLink, Link } from 'react-router';
import { connect } from 'react-redux';
import { Dropdown, Header, Icon, Menu } from 'semantic-ui-react';

import AuthThunks from '../../modules/auth/auth.thunk';
import { isRoleAuthorized } from '../../modules/auth/auth.wrappers';

// Style
import './navBar.component.scss';

// NavBar Component
class NavBarComponent extends React.Component {

  isAuthorized = (child, Roles) => {
    const authorized = this.props.auth.isAuthenticated && isRoleAuthorized(Roles, this.props.auth.user.role);
    return authorized && child;
  }

  isActiveURL = (url) => {
    return this.props.location && this.props.location.pathname && this.props.location.pathname.startsWith(url);
  }

  renderDropdown = (loading) => {
    const item = [];
    item.push(<Icon key='icon' name={loading ? 'circle notched' : 'user'} loading={loading} size='large'/>);
    item.push(this.props.auth.user.displayName);
    return item;
  }

  render = () => {
    const { logout, exportDocktor, isExportFetching } = this.props;
    const isAuthorized = this.isAuthorized;
    return (
      <Menu inverted className='navbar'>
        <Menu.Item  as={IndexLink} to='/' header color='blue' active>
          <Icon name='dashboard' size='large'/>
          <Header.Content>
            D.A.D
          </Header.Content>
        </Menu.Item>
        {isAuthorized(<Menu.Item active={this.isActiveURL('/users')} as={Link} to='/users'>Users</Menu.Item>)}
        {isAuthorized(
          <Menu.Menu position='right'>
            <Menu.Item as={Dropdown} trigger={this.renderDropdown(isExportFetching)}>
              <Dropdown.Menu>
                {isAuthorized(
                  <Dropdown.Item onClick={exportDocktor} disabled={isExportFetching}><Icon name='download' />Export</Dropdown.Item>,
                )}
                <Dropdown.Item as={Link} to='/settings'><Icon name='settings' />Settings</Dropdown.Item>
                <Dropdown.Item onClick={logout} ><Icon name='sign out' />Logout</Dropdown.Item>
              </Dropdown.Menu>
            </Menu.Item>
          </Menu.Menu>
        )}
      </Menu>
    );
  }
}

NavBarComponent.propTypes = {
  location: React.PropTypes.object.isRequired,
  auth: React.PropTypes.object.isRequired,
  logout: React.PropTypes.func.isRequired,
  exportDocktor: React.PropTypes.func.isRequired,
  isExportFetching: React.PropTypes.bool
};

// Function to map state to container props
const mapStateToProps = (state) => {
  return {
    auth: state.auth,
    location: state.routing.locationBeforeTransitions,
    //isExportFetching: state.export.isFetching
  };
};

// Function to map dispatch to container props
const mapDispatchToProps = (dispatch) => {
  return {
    logout: () => {
      dispatch(AuthThunks.logoutUser());
    },
    exportDocktor: () => {
      dispatch(ExportThunks.exportAll());
    }
  };
};

// Redux container to Sites component
const NavBar = connect(
  mapStateToProps,
  mapDispatchToProps
)(NavBarComponent);

export default NavBar;
