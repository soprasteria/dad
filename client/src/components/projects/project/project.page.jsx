// React
import React from 'react';
import { Link } from 'react-router';
import { connect } from 'react-redux';
import { push } from 'react-router-redux';
import { Button, Container, Divider, Form, Header, Icon, Message, Table, Segment } from 'semantic-ui-react';

import Joi from 'joi-browser';

import Matrix from './matrix/matrix.component';
import Box from '../../common/box.component';

// Thunks / Actions
import ProjectsThunks from '../../../modules/projects/projects.thunks';
import EntitiesThunks from '../../../modules/entities/entities.thunks';
import ServicesThunks from '../../../modules/services/services.thunks';
import UsersThunks from '../../../modules/users/users.thunks';

import { getEntitiesAsOptions, getByType } from '../../../modules/entities/entities.selectors';
import { groupByPackage } from '../../../modules/services/services.selectors';
import { getUsersAsOptions } from '../../../modules/users/users.selectors';

import { parseError } from '../../../modules/utils/forms';

import { AUTH_CP_ROLE } from '../../../modules/auth/auth.constants';

// Style
import './project.page.scss';

// Project Component
class ProjectComponent extends React.Component {

  state = { errors: { details: [], fields: {} }, project: {}, matrix: {} }

  schema = Joi.object().keys({
    name: Joi.string().trim().required().label('Project Name'),
    domain: Joi.string().trim().required().label('Domain'),
    projectManager: Joi.string().trim().alphanum().empty('').label('Project Manager'),
    serviceCenter: Joi.string().trim().alphanum().empty('').label('Service Center'),
    businessUnit: Joi.string().trim().alphanum().empty('').label('Business Unit')
  }).or('serviceCenter', 'businessUnit').label('Service Center or Business Unit');

  componentWillMount = () => {
    const matrix = {};
    const project = this.props.project || { matrix: [] };
    project.matrix.forEach((m) => matrix[m.service] = m);
    this.setState({ project: { ...project }, errors: { details: [], fields:{} }, matrix });
  }

  componentWillReceiveProps = (nextProps) => {
    const matrix = {};
    const project = nextProps.project || { matrix: [] };
    project.matrix.forEach((m) => matrix[m.service] = m);
    this.setState({ project: { ...project }, errors: { details: [], fields:{} }, matrix });
  }

  componentDidMount = () => {
    const { projectId } = this.props;
    Promise.all([
      this.props.fetchEntities(),
      this.props.fetchServices(),
      this.props.fetchUsers()
    ]).then(()=>{
      if (projectId) {
        this.props.fetchProject(projectId);
      }
    });
    if (!projectId) {
      window.scrollTo(0, 0);
    }
  }

  componentDidUpdate = (prevProps) => {
    if (prevProps.isFetching) {
      window.scrollTo(0, 0);
    }
  }

  handleChange = (e, { name, value }) => {
    const { project, errors } = this.state;
    const state = {
      project: { ...project, [name]: value },
      errors: { details: [...errors.details], fields: { ...errors.fields } }
    };
    delete state.errors.fields[name];
    this.setState(state);
  }

  handleMatrix = (id, newMatrix) => {
    const { matrix } = this.state;
    newMatrix.service = id;
    this.setState({ matrix: { ...matrix, [id]: newMatrix } });
  }

  isFormValid = () => {
    const { error } = Joi.validate(this.state.project, this.schema, { abortEarly: false, allowUnknown: true });
    if (error) {
      const errors = parseError(error);
      if (errors.fields['Service Center or Business Unit']) {
        errors.fields.serviceCenter = true;
        errors.fields.businessUnit = true;
        delete errors.fields['Service Center or Business Unit'];
      }
      window.scrollTo(0, 0);
      this.setState({ errors:  errors });
    }
    return !Boolean(error);
  }

  handleSubmit = (e) => {
    e.preventDefault();
    if (this.isFormValid()) {
      const { project, matrix } = this.state;
      const modifiedProject = { ...project, matrix:Object.values(matrix) };
      this.props.onSave(modifiedProject);
    }
  }

  renderServices = (project, services, isServicesFetching) => {
    if (isServicesFetching) {
      return <div />;
    }
    const readOnly = this.isReadonly();
    return Object.entries(services).map(([pckg, servicesList]) => {
      return (
        <Table key={pckg} celled striped compact>
          <Table.Header>
            <Table.Row>
              <Table.HeaderCell width='six'>{pckg}</Table.HeaderCell>
              <Table.HeaderCell width='two'>Progress</Table.HeaderCell>
              <Table.HeaderCell width='two'>Goal</Table.HeaderCell>
              <Table.HeaderCell width='six'>Comment</Table.HeaderCell>
            </Table.Row>
          </Table.Header>
          <Table.Body>
            {servicesList.map(service => {
              return <Matrix readOnly={readOnly} serviceId={service.id} key={service.id} matrix={this.state.matrix[service.id] || {}} service={service} onChange={this.handleMatrix}/>;
            })}
          </Table.Body>
        </Table>
      );
    });
  }

  renderDropdown = (name, label, value, placeholder, width, options, isFetching, errors ) => {
    if (this.isReadonly()) {
      const option = options.find(elm => elm.value === value);
      return (
        <Form.Input readOnly label={label} value={(option && option.text) || ''} onChange={this.handleChange}
          type='text' autoComplete='off' placeholder={`No ${label}`} width={width}
        />
      );
    }
    return (
      <Form.Dropdown placeholder={placeholder} fluid search selection loading={isFetching}  width={width}
        label={label} name={name} options={options} value={value || ''} onChange={this.handleChange} error={errors.fields[name]}
      />
    );
  }

  isReadonly = () => this.props.auth.user.role === AUTH_CP_ROLE

  render = () => {
    const { isFetching, serviceCenters, businessUnits, isEntitiesFetching, services, isServicesFetching, users, projectId } = this.props;
    const { project, errors } = this.state;
    const readOnly = this.isReadonly();
    return (
      <Container className='project-page'>
        <Segment loading={isFetching || isServicesFetching} padded>
          <Header as='h1'>
            <Link to={'/projects'}>
              <Icon name='arrow left' fitted/>
            </Link>
            {projectId ? project.name : 'New Project'}
            {project.url && <Button as='a' href={project.url} content='URL' icon='linkify' labelPosition='left' color='blue' floated='right' />}
          </Header>
          <Divider hidden/>
          <Box icon='settings' title='Details' stacked={Boolean(projectId)}>
            <Form error={Boolean(errors.details.length)}>
              <Form.Group>
                <Form.Input readOnly={readOnly} label='Name' value={project.name || ''} onChange={this.handleChange}
                  type='text' name='name' autoComplete='off' placeholder='Project Name' width='four' error={errors.fields['name']}
                />
                <Form.Input readOnly={readOnly} label='Domain' value={project.domain || ''} onChange={this.handleChange}
                    type='text' name='domain' autoComplete='off' placeholder='Project Domain' width='four' error={errors.fields['domain']}
                />
                {this.renderDropdown('projectManager', 'Project Manager', project.projectManager, 'Select Project Manager...', 'eight', users, isEntitiesFetching, errors)}
              </Form.Group>

              <Form.Group widths='two'>
                {this.renderDropdown('serviceCenter', 'Service Center', project.serviceCenter, 'Select Service Center...', 'eight', serviceCenters, isEntitiesFetching, errors)}
                {this.renderDropdown('businessUnit', 'Business Unit', project.businessUnit, 'Select Business Unit...', 'eight', businessUnits, isEntitiesFetching, errors)}
              </Form.Group>
              <Message error list={errors.details}/>
            </Form>
          </Box>
          <Divider hidden/>
          {this.renderServices(project, services, isServicesFetching)}
          <Divider hidden/>
          <Button fluid color='green' onClick={this.handleSubmit}>Save</Button>
        </Segment>
      </Container>
    );
  }
}

ProjectComponent.propTypes = {
  auth: React.PropTypes.object.isRequired,
  project: React.PropTypes.object,
  isFetching: React.PropTypes.bool,
  businessUnits: React.PropTypes.array,
  serviceCenters: React.PropTypes.array,
  isEntitiesFetching: React.PropTypes.bool,
  users: React.PropTypes.array,
  services: React.PropTypes.object,
  isServicesFetching: React.PropTypes.bool,
  projectId: React.PropTypes.string,
  fetchProject: React.PropTypes.func.isRequired,
  fetchEntities: React.PropTypes.func.isRequired,
  fetchServices: React.PropTypes.func.isRequired,
  fetchUsers: React.PropTypes.func.isRequired,
  onSave: React.PropTypes.func.isRequired
};

const mapStateToProps = (state, ownProps) => {
  const paramId = ownProps.params.id;
  const projects = state.projects;
  const project = projects.selected;
  const emptyProject = { matrix: [] };
  const isFetching = paramId && (paramId !== project.id || project.isFetching);
  const entities = Object.values(state.entities.items);
  const services = groupByPackage(state.services.items);
  const isServicesFetching = state.services.isFetching;
  const users = Object.values(state.users.items);
  return {
    auth: state.auth,
    project: { ...emptyProject, ...projects.items[paramId] },
    isFetching,
    projectId: paramId,
    businessUnits: getEntitiesAsOptions(getByType(entities, 'businessUnit')),
    serviceCenters: getEntitiesAsOptions(getByType(entities, 'serviceCenter')),
    users: getUsersAsOptions(users),
    isEntitiesFetching: state.entities.isFetching,
    services,
    isServicesFetching
  };
};

const mapDispatchToProps = dispatch => ({
  fetchProject: id => dispatch(ProjectsThunks.fetch(id)),
  fetchEntities: () => dispatch(EntitiesThunks.fetchIfNeeded()),
  fetchServices: () => dispatch(ServicesThunks.fetchIfNeeded()),
  fetchUsers: () => dispatch(UsersThunks.fetchIfNeeded()),
  onSave: project => dispatch(ProjectsThunks.save(project, push('/projects')))
});

const ProjectPage = connect(
  mapStateToProps,
  mapDispatchToProps
)(ProjectComponent);

export default ProjectPage;
