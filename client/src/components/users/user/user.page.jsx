// React
import React from 'react';
import { Link } from 'react-router';
import { connect } from 'react-redux';
import { push } from 'react-router-redux';
import { Button, Container, Divider, Dropdown, Form, Header, Icon, Label, Segment } from 'semantic-ui-react';
import { ALL_ROLES, getRoleColor, getRoleIcon, getRoleLabel } from '../../../modules/auth/auth.constants';

// Thunks / Actions
import UsersThunks from '../../../modules/users/users.thunks';
import OrganizationsThunks from '../../../modules/organizations/organizations.thunks';

import { getOrganizationsAsOptions } from '../../../modules/organizations/organizations.selectors';

// Style
import './user.page.scss';

// User Component
class UserComponent extends React.Component {

  state = { user: {} }

  componentWillMount = () => {
    this.setState({ user: { ...this.props.user } });
  }

  componentWillReceiveProps = (nextProps) => {
    this.setState({ user: { ...nextProps.user } });
  }

  componentDidMount = () => {
    const { userId } = this.props;
    Promise.all([this.props.fetchOrganizations()]).then(()=>{
      this.props.fetchUser(userId);
    });
  }

  handleChange = (e, { name, value, checked }) => {
    const { user } = this.state;
    const state = {
      user: { ...user, [name]:value || checked }
    };
    this.setState(state);
  }

  handleSubmit = (e) => {
    e.preventDefault();
    const stateUser = this.state.user;
    const user = { ...stateUser, projects: [...stateUser.projects], organizations: [...stateUser.organizations] };
    this.props.onSave(user);
  }

  renderDropDownButton = (user, isFetching) => {
    return (
      <Button loading={isFetching} color={getRoleColor(user.role)} className='role' onClick={e => e.preventDefault()}>
        <Icon name={getRoleIcon(user.role)} />
        {getRoleLabel(user.role)}
      </Button>
    );
  }

  renderRoleDropdown = (user, isFetching) => {
    const options = ALL_ROLES.map(role => {
      return {
        icon: <Icon name={getRoleIcon(role)} color={getRoleColor(role)} />,
        value: role,
        text: getRoleLabel(role)
      };
    });

    return (
      <Dropdown pointing trigger={this.renderDropDownButton(user, isFetching)} onChange={this.handleChange} options={options}
        icon={null} value={user.role} name='role'
      />
    );
  }

  render = () => {
    const { isFetching, projects, organizations, isOrganizationsFetching } = this.props;
    const { user } = this.state;
    return (
      <Container className='user-page'>
        <Segment loading={isFetching} padded>
          <Header as='h1'>
            <Link to={'/users'}>
              <Icon name='arrow left' fitted/>
            </Link>
            {`${user.displayName} (${user.username})`}
          </Header>
          <Divider hidden/>
          <Form onSubmit={this.handleSubmit}>
            <Form.Group widths='two' >
              <Form.Input readOnly label='Username' value={user.username || ''} onChange={this.handleChange}
                type='text' name='username' autoComplete='off' placeholder='Username'
              />
              <Form.Input readOnly label='Email Address' value={user.email || ''} onChange={this.handleChange}
                  type='text' name='email' autoComplete='off' placeholder='A valid email address'
              />
            </Form.Group>

            <Form.Group widths='two'>
              <Form.Input readOnly label='First Name' value={user.firstName || ''} onChange={this.handleChange}
                type='text' name='firstName' autoComplete='off' placeholder='First Name'
              />
              <Form.Input readOnly label='Last Name' value={user.lastName || ''} onChange={this.handleChange}
                  type='text' name='lastName' autoComplete='off' placeholder='Last Name'
              />
            </Form.Group>

            <Divider hidden/>

            <Form.Group>
              <Form.Field width='two'>
                <Label size='large' className='form-label' content='Role' />
              </Form.Field>
              <Form.Field width='fourteen'>
                {this.renderRoleDropdown(user, isFetching)}
              </Form.Field>
            </Form.Group>

            <Form.Group>
              <Form.Field width='two'>
                <Label size='large' className='form-label' content='Projects' />
              </Form.Field>
              <Form.Dropdown width='fourteen' placeholder='Select projects' fluid multiple search selection
                name='projects' options={projects} value={user.projects || []} onChange={this.handleChange}
              />
            </Form.Group>

            <Form.Group>
              <Form.Field width='two'>
                <Label size='large' className='form-label' content='Entities' />
              </Form.Field>
              <Form.Dropdown width='fourteen' placeholder='Select organizations' fluid multiple search selection loading={isOrganizationsFetching}
                name='organizations' options={organizations} value={user.organizations || []} onChange={this.handleChange}
              />
            </Form.Group>

            <Divider hidden/>

            <Button fluid color='green' content='Save' loading={isFetching} />
          </Form>
        </Segment>
      </Container>
    );
  }
}

UserComponent.propTypes = {
  user: React.PropTypes.object,
  isFetching: React.PropTypes.bool,
  projects: React.PropTypes.array,
  organizations: React.PropTypes.array,
  isOrganizationsFetching: React.PropTypes.bool,
  userId: React.PropTypes.string.isRequired,
  fetchUser: React.PropTypes.func.isRequired,
  fetchOrganizations: React.PropTypes.func.isRequired,
  onSave: React.PropTypes.func.isRequired
};

const mapStateToProps = (state, ownProps) => {
  const paramId = ownProps.params.id;
  const users = state.users;
  const user = users.selected;
  const emptyUser = { projects: [], organizations: [] };
  const isFetching = paramId && (paramId !== user.id || user.isFetching);
  return {
    user: users.items[user.id] || emptyUser,
    isFetching,
    userId: paramId,
    projects: [{ text: 'Test project', value: '587fe317693c38712351b7cb' } ],
    organizations: getOrganizationsAsOptions(state.organizations.items),
    isOrganizationsFetching: state.organizations.isFetching
  };
};

const mapDispatchToProps = dispatch => ({
  fetchUser: id => dispatch(UsersThunks.fetch(id)),
  fetchOrganizations : () => dispatch(OrganizationsThunks.fetchIfNeeded()),
  onSave: user => dispatch(UsersThunks.save(user, push('/users')))
});

const UserPage = connect(
  mapStateToProps,
  mapDispatchToProps
)(UserComponent);

export default UserPage;
