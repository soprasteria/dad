// React
import React from 'react';
import { connect } from 'react-redux';
import { Card, Container, Icon, Input, Label, Segment } from 'semantic-ui-react';
import DebounceInput from 'react-debounce-input';

// API Fetching
import UsersThunks from '../../modules/users/users.thunks';
import UsersActions from '../../modules/users/users.actions';

// Selectors
import { getFilteredUsers } from '../../modules/users/users.selectors';

// Components
import UserCard from './user/user.card.component';

import './users.page.scss';

//Site Component using react-leaflet
class Users extends React.Component {

  componentWillMount = () => {
    document.title = 'D.A.D - Users';
    this.props.fetchUsers();
  }

  renderCards = (users) => {
    if (users.length) {
      return (
        <Card.Group className='centered'>
          {users.map(user => {
            return (
              <UserCard user={user} key={user.id} />
            );
          })}
        </Card.Group>
      );
    }
    return <p>No users found.</p>;
  }

  render = () => {
    const { users, filterValue, isFetching, changeFilter } = this.props;
    return (
      <Container fluid className='users-page'>
        <Segment.Group raised>
          <Segment>
            <Input fluid icon labelPosition='left corner'>
              <Label corner='left' icon='search' />
              <DebounceInput
                placeholder='Search...'
                minLength={1}
                debounceTimeout={300}
                onChange={(event) => changeFilter(event.target.value)}
                value={filterValue}
              />
              <Icon link name='remove' onClick={() => changeFilter('')}/>
            </Input>
          </Segment>
          <Segment loading={isFetching}>
            {this.renderCards(users)}
          </Segment>
        </Segment.Group>
      </Container>
    );
  }
}

Users.propTypes = {
  users: React.PropTypes.array,
  filterValue: React.PropTypes.string,
  isFetching: React.PropTypes.bool,
  fetchUsers: React.PropTypes.func.isRequired,
  changeFilter: React.PropTypes.func.isRequired
};

// Function to map state to container props
const mapStateToUsersProps = (state) => {
  const filterValue = state.users.filterValue;
  const users = getFilteredUsers(state.users.items, filterValue);
  const isFetching = state.users.isFetching;
  return { filterValue, users, isFetching };
};

// Function to map dispatch to container props
const mapDispatchToUsersProps = (dispatch) => {
  return {
    fetchUsers : () => dispatch(UsersThunks.fetchIfNeeded()),
    changeFilter: filterValue => dispatch(UsersActions.changeFilter(filterValue))
  };
};

// Redux container to Sites component
const UsersPage = connect(
  mapStateToUsersProps,
  mapDispatchToUsersProps
)(Users);

export default UsersPage;
