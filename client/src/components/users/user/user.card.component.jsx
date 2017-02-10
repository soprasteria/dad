// React
import React from 'react';
import { connect } from 'react-redux';
import { Link } from 'react-router';
import { Button, Card, Dropdown, Icon } from 'semantic-ui-react';

import { ALL_ROLES, AUTH_ADMIN_ROLE, getRoleColor, getRoleIcon, getRoleLabel } from '../../../modules/auth/auth.constants';
import UsersThunks from '../../../modules/users/users.thunks';

import './user.card.component.scss';

// UserCard Component
class UserCardComponent extends React.Component {

  handleChange = (e, { value }) => {
    const oldUser = this.props.user;
    const userToSave = {
      ...oldUser,
      Role: value
    };
    this.props.saveUser(userToSave);
  }

  renderDropDown = (user) => {
    return (
      <Button loading={user.isFetching} color={getRoleColor(user.role)} compact size='small'>
        <Icon name={getRoleIcon(user.role)}  />
        {getRoleLabel(user.role)}
      </Button>
    );
  }

  render = () => {
    const user = this.props.user;
    const connectedUser = this.props.auth.user;
    const disabled = connectedUser.role !== AUTH_ADMIN_ROLE;
    const options = ALL_ROLES.map(role => {
      return { icon: <Icon name={getRoleIcon(role)} color={getRoleColor(role) || null} />, value: role, text: getRoleLabel(role) };
    });

    return (
      <Card className='user-card' raised>
        <Card.Content>
          <Link to={`/users/${user.id}`} title={user.displayName}>
            {`${user.lastName.toUpperCase()} ${user.firstName}`}
          </Link>
          <Dropdown trigger={this.renderDropDown(user)} compact onChange={this.handleChange} options={options}
            icon={null} button disabled={disabled} value={user.role} pointing='right' className='tiny attached'
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
