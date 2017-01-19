// React
import React from 'react';
import { connect } from 'react-redux';
import { Link } from 'react-router';
import { Button, Card, Dropdown, Icon } from 'semantic-ui-react';

import { ALL_ROLES, AUTH_ADMIN_ROLE, AUTH_RI_ROLE, getRoleColor, getRoleIcon, getRoleLabel } from '../../../modules/auth/auth.constants';
import UsersThunks from '../../../modules/users/users.thunks';

import './user.card.component.scss';

// UserCard Component
class UserCardComponent extends React.Component {

  state = { isFetching: false }

  componentWillReceiveProps = () => {
    this.setState({ isFetching: false });
  }

  handleChange = (e, { value }) => {
    const oldUser = this.props.user;
    const userToSave = {
      ...oldUser,
      Role: value
    };
    this.setState({ isFetching : true });
    this.props.saveUser(userToSave);
  }

  renderDropDown = (user, isFetching) => {
    return (
      <Button loading={isFetching} color={getRoleColor(user.role)} compact size='small'>
        <Icon name={getRoleIcon(user.role)}  />
        {getRoleLabel(user.role)}
      </Button>
    );
  }

  render = () => {
    const { isFetching } = this.state;
    const user = this.props.user;
    const connectedUser = this.props.auth.user;
    const disabled = connectedUser.role !== AUTH_ADMIN_ROLE;
    const options = ALL_ROLES.map(role => {
      return { icon: <Icon name={getRoleIcon(role)} color={getRoleColor(role) || null} />, value: role, text: getRoleLabel(role) };
    });
    const canGoToProfile = connectedUser.role === AUTH_ADMIN_ROLE || connectedUser.role === AUTH_RI_ROLE;
    return (
      <Card className='user-card'>
        <Card.Content>
          {
            canGoToProfile ?
              <Link to={`/users/${user.id}`}>
                {user.displayName}
              </Link>
            :
              user.displayName
          }
          <Dropdown trigger={this.renderDropDown(user, isFetching)} compact onChange={this.handleChange} options={options}
            icon={null} button disabled={disabled} value={user.role} pointing className='tiny top right attached'
          />
        </Card.Content>
        <Card.Content extra>
          <div className='email' title={user.email}>
            <a href={`mailto:${user.email}`}><Icon name='mail' />{user.email}</a>
          </div>
        </Card.Content>
      </Card>
    );
  }
}

UserCardComponent.propTypes = {
  user: React.PropTypes.object,
  auth: React.PropTypes.object,
  saveUser: React.PropTypes.func.isRequired
};

// Function to map state to container props
const mapStateToProps = (state) => {
  return {
    auth: state.auth,
  };
};

// Function to map dispatch to container props
const mapDispatchToProps = (dispatch) => {
  return {
    saveUser: (user) => {
      dispatch(UsersThunks.save(user));
    }
  };
};

// Redux container to Sites component
const UserCard = connect(
  mapStateToProps,
  mapDispatchToProps
)(UserCardComponent);

export default UserCard;
