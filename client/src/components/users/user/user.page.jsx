// React
import React from 'react';
import { Link } from 'react-router';
import { connect } from 'react-redux';
import { push } from 'react-router-redux';
import { Form, Button, Segment, Label, Icon, Dropdown } from 'semantic-ui-react';
import { ALL_ROLES, getRoleLabel, getRoleColor, getRoleIcon } from '../../../modules/auth/auth.constants';

// Thunks / Actions
import UsersThunks from '../../../modules/users/users.thunks';

// Style
import './user.page.scss';

// User Component
class UserComponent extends React.Component {

  state = { errors: { details: [], fields: {} }, user: {} }

  schema = Joi.object().keys({
    projects: Joi.array().items(Joi.string().alphanum().trim()),
    entities: Joi.array().items(Joi.string().alphanum().trim()),
  })

  componentWillMount = () => {
    this.setState({ user: { ...this.props.user }, errors: { details: [], fields:{} } });
  }

  componentWillReceiveProps = (nextProps) => {
    this.setState({ user: { ...nextProps.user }, errors: { details: [], fields:{} } });
  }

  componentDidMount = () => {
    const { userId } = this.props;
    this.props.fetchUser(userId);
  }

  handleChange = (e, { name, value, checked }) => {
    const { user, errors } = this.state;
    const state = {
      user: { ...user, [name]:value || checked },
      errors: { details: [...errors.details], fields: { ...errors.fields } }
    };
    delete state.errors.fields[name];
    this.setState(state);
  }

  isFormValid = () => {
    const { error } = Joi.validate(this.state.user, this.schema, { abortEarly: false, allowUnknown: true });
    error && this.setState({ errors: parseError(error) });
    return !Boolean(error);
  }

  onSave = (e) => {
    e.preventDefault();
    if (this.isFormValid()) {
      const stateUser = this.state.user;
      const user = { ...stateUser, projects: [...stateUser.projects], entities: [...stateuser.entities] };
      this.props.onSave(user);
    }
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
      <Dropdown trigger={this.renderDropDownButton(user, isFetching)} onChange={this.handleChange} options={options}
        icon={null} value={user.role} name='role'
      />
    );
  }

  render = () => {
    const { isFetching } = this.props;
    const { user } = this.state;
    return (
      <Segment loading={isFetching}>
              <h1>
                <Link to={'/users'}>
                  <Icon name='arrow left' fitted/>
                </Link>
                {`${user.displayName} (${user.username})`}
              </h1>

              <Form className='user-form'>
                <Form.Group>
                  <Form.Field width='two'>
                    <Label size='large' content='Role' />
                  </Form.Field>

                  <Form.Field width='fourteen'>
                    {this.renderRoleDropdown(user, isFetching)}
                  </Form.Field>
                </Form.Group>

                <Form.Group widths='two'>
                  <Form.Input required readOnly label='Username' value={user.username || ''} onChange={this.handleChange}
                    type='text' name='username' autoComplete='off' placeholder='Username'
                  />
                  <Form.Input required readOnly label='Email Address' value={user.email || ''} onChange={this.handleChange}
                      type='text' name='email' autoComplete='off' placeholder='A valid email address'
                  />
                </Form.Group>

                <Form.Group widths='two'>
                  <Form.Input required readOnly label='First Name' value={user.firstName || ''} onChange={this.handleChange}
                    type='text' name='firstName' autoComplete='off' placeholder='First Name'
                  />
                  <Form.Input required readOnly label='Last Name' value={user.lastName || ''} onChange={this.handleChange}
                      type='text' name='lastName' autoComplete='off' placeholder='Last Name'
                  />
                </Form.Group>
                <Button fluid onClick={this.onSave}>Save</Button>
              </Form>
      </Segment>
    );
  }
}

UserComponent.propTypes = {
  user: React.PropTypes.object,
  isFetching: React.PropTypes.bool,
  userId: React.PropTypes.string.isRequired,
  tags: React.PropTypes.object,
  fetchUser: React.PropTypes.func.isRequired,
  fetchTags: React.PropTypes.func.isRequired,
  onSave: React.PropTypes.func.isRequired
};

const mapStateToProps = (state, ownProps) => {
  const paramId = ownProps.params.id;
  const users = state.users;
  const user = users.selected;
  const emptyUser = { tags: [] };
  const isFetching = paramId && (paramId !== user.id);
  return {
    user: users.items[user.id] || emptyUser,
    isFetching,
    userId: paramId,
    tags: state.tags
  };
};

const mapDispatchToProps = dispatch => ({
  fetchUser: id => dispatch(UsersThunks.fetch(id)),
  fetchTags: () => dispatch(TagsThunks.fetchIfNeeded()),
  onSave: user => dispatch(UsersThunks.save(user, push('/users')))
});

const UserPage = connect(
  mapStateToProps,
  mapDispatchToProps
)(UserComponent);

export default UserPage;
