// React
import React from 'react';
import PropTypes from 'prop-types';
import { IndexLink, Link } from 'react-router';
import { connect } from 'react-redux';
import { Dropdown, Header, Icon, Menu } from 'semantic-ui-react';

import AuthThunks from '../../modules/auth/auth.thunk';
import ExportThunks from '../../modules/export/export.thunk';
import LanguagesThunks from '../../modules/languages/languages.thunks';
import { isRoleAuthorized } from '../../modules/auth/auth.wrappers';
import { flattenLanguages } from '../../modules/languages/languages.selectors';

// Style
import './navBar.component.scss';

// NavBar Component
class NavBarComponent extends React.Component {

  componentDidMount = () => {
    this.props.fetchLanguages();
  }

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
    const name = `${this.props.auth.user.lastName ? this.props.auth.user.lastName.toUpperCase() : ''} ${this.props.auth.user.firstName}`;
    item.push(name);
    return item;
  }

  render = () => {
    const { logout, exportData, languages, language, selectLanguage, isExportFetching } = this.props;
    const isAuthorized = this.isAuthorized;
    const userId = this.props.auth.user.id;
    return (
      <Menu inverted className='navbar'>
        <Menu.Item  as={IndexLink} to='/' header color='blue' active>
          <Icon name='dashboard' size='large'/>
          <Header.Content>
            D.A.D
          </Header.Content>
        </Menu.Item>
        {isAuthorized(<Menu.Item active={this.isActiveURL('/projects')} as={Link} to='/projects'>Projects</Menu.Item>)}
        {isAuthorized(<Menu.Item active={this.isActiveURL('/users')} as={Link} to='/users'>Users</Menu.Item>)}
        {isAuthorized(
          <Menu.Menu position='right'>
            <Menu.Item as={Dropdown} text={language}>
              <Dropdown.Menu>
                {languages.map((language, index) => (
                  <Dropdown.Item key={index} onClick={selectLanguage.bind(this, language)}>{language}</Dropdown.Item>
                ))}
              </Dropdown.Menu>
            </Menu.Item>
            <Menu.Item as={Dropdown} trigger={this.renderDropdown(isExportFetching)}>
              <Dropdown.Menu>
                {isAuthorized(
                  <Dropdown.Item onClick={exportData.bind(this, language)} disabled={isExportFetching}><Icon name='download' />Export</Dropdown.Item>,
                )}
                <Dropdown.Item as={Link} to={`/users/${userId}`} ><Icon name='settings' />Profile</Dropdown.Item>
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
  location: PropTypes.object.isRequired,
  auth: PropTypes.object.isRequired,
  logout: PropTypes.func.isRequired,
  exportData: PropTypes.func.isRequired,
  isExportFetching: PropTypes.bool,
  fetchLanguages: PropTypes.func.isRequired,
  languages: PropTypes.array.isRequired,
  language: PropTypes.string.isRequired,
  selectLanguage: PropTypes.func.isRequired
};

// Function to map state to container props
const mapStateToProps = (state) => {
  const languages = state.languages;
  return {
    auth: state.auth,
    location: state.routing.locationBeforeTransitions,
    languages: flattenLanguages(languages.items),
    language: languages.language
    //isExportFetching: state.export.isFetching
  };
};

// Function to map dispatch to container props
const mapDispatchToProps = (dispatch) => {
  return {
    logout: () => {
      dispatch(AuthThunks.logoutUser());
    },
    selectLanguage: (language) => {
      dispatch(LanguagesThunks.selectLanguage(language));
    },
    exportData: (language) => {
      dispatch(ExportThunks.exportAll(language));
    },
    fetchLanguages: () => {
      dispatch(LanguagesThunks.fetchAll());
    }
  };
};

// Redux container to Sites component
const NavBar = connect(
  mapStateToProps,
  mapDispatchToProps
)(NavBarComponent);

export default NavBar;
