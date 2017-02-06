// React
import React from 'react';
import { Link } from 'react-router';
import { connect } from 'react-redux';
import { push } from 'react-router-redux';
import { Button, Container, Divider, Form, Icon, Label, Message, Table, Segment } from 'semantic-ui-react';

import Joi from 'joi-browser';

import Matrix from './matrix/matrix.component';
import Box from '../../common/box.component';

// Thunks / Actions
import ProjectsThunks from '../../../modules/projects/projects.thunks';
import EntitiesThunks from '../../../modules/entities/entities.thunks';
import ServicesThunks from '../../../modules/services/services.thunks';
import UsersThunks from '../../../modules/users/users.thunks';
import ProjectsActions from '../../../modules/projects/projects.actions';
import ModalActions from '../../../modules/modal/modal.actions';
import ToastsActions from '../../../modules/toasts/toasts.actions';

import { getEntitiesAsOptions, getByType } from '../../../modules/entities/entities.selectors';
import { groupByPackage } from '../../../modules/services/services.selectors';
import { getUsersAsOptions } from '../../../modules/users/users.selectors';

import { parseError } from '../../../modules/utils/forms';

import { AUTH_CP_ROLE, AUTH_RI_ROLE, AUTH_ADMIN_ROLE } from '../../../modules/auth/auth.constants';

// Style
import './project.page.scss';

// Project Component
class ProjectComponent extends React.Component {

  state = { errors: { details: [], fields: {} }, project: {}, matrix: {} }

  schema = Joi.object().keys({
    name: Joi.string().trim().required().label('Project Name'),
    domain: Joi.string().trim().empty('').label('Domain'),
    projectManager: Joi.string().trim().alphanum().empty('').label('Project Manager'),
    serviceCenter: Joi.string().trim().alphanum().empty('').label('Service Center'),
    businessUnit: Joi.string().trim().alphanum().empty('').label('Business Unit')
  }).or('serviceCenter', 'businessUnit').label('Service Center or Business Unit');

  componentWillMount = () => {
    const matrix = {};
    const project = this.props.project;
    project.matrix.forEach((m) => matrix[m.service] = m);
    this.setState({ project: { ...project }, errors: { details: [], fields:{} }, matrix });
  }

  componentWillReceiveProps = (nextProps) => {
    const project = nextProps.project;
    if(!project.isEditing) {
      const matrix = {};
      project.matrix.forEach((m) => matrix[m.service] = m);
      this.setState({ project: { ...project }, errors: { details: [], fields:{} }, matrix });
    } else {
      this.setState({ project: { ...this.state.project, urls: [...project.urls] } });
    }
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
      this.refs.details.setState({ stacked:false });
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

  handleRemove = (e) => {
    e.preventDefault();
    this.props.onDelete(this.state.project);
  }

  renderServices = (project, services, isFetching, readOnly) => {
    if (isFetching) {
      return <p>Fetching Matrix...</p>;
    }
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

  renderDropdown = (name, label, value, placeholder, width, options, isFetching, errors, readOnly) => {
    if (readOnly) {
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

  render = () => {
    const {
      isFetching, serviceCenters, businessUnits,
      isEntitiesFetching, services, isServicesFetching,
      users, projectId, onCreateUrl, onEditUrl, onRemoveUrl,
      canEditDetails
    } = this.props;
    const { project, errors } = this.state;
    const fetching = isFetching || isServicesFetching;
    const authUser = this.props.auth.user;
    const canEditMatrix = canEditDetails || (authUser.role === AUTH_CP_ROLE && project.projectManager === authUser.id);
    const createUrl = (e) => {
      e.preventDefault();
      onCreateUrl(projectId);
    };
    const editUrl = (index, url) => e => {
      e.preventDefault();
      onEditUrl(projectId, index, url);
    };
    const removeUrl = (index) => e => {
      e.preventDefault();
      onRemoveUrl(projectId, index);
    };
    return (
      <Container className='project-page'>

        <Segment loading={fetching} padded>

          <Form>
          <h1 className='layout horizontal center justified'>
            <Link to={'/projects'}>
              <Icon name='arrow left' fitted/>
            </Link>
            <Form.Input className='flex projectName' readOnly={!canEditDetails} value={project.name || ''} onChange={this.handleChange}
              type='text' name='name' autoComplete='off' placeholder='Project Name' error={errors.fields['name']}
            />
            {(!isFetching && canEditDetails) && <Button color='red' icon='trash' labelPosition='left' title='Delete project' content='Delete Project' onClick={this.handleRemove} />}
          </h1>

          <Divider hidden/>
            <Form.TextArea readOnly={!canEditMatrix} label='Description' value={project.description || ''} onChange={this.handleChange} autoHeight
                  type='text' name='description' autoComplete='off' placeholder='Project description' width='sixteen' error={errors.fields['description']}/>
            <Form.Group>
              <Form.Field width='two'>
                <Label size='large' className='form-label' content='URLs' />
              </Form.Field>
              <Form.Field width='fourteen'>
                <Label.Group>
                  {project.urls && project.urls.map((url, index) => {
                    return (
                      <Label as='a' href={url.link} color='blue' key={index} image>
                        <Icon name='linkify' />
                        {url.name}
                        {canEditDetails &&
                          <Label.Detail>
                            <Icon link fitted name='edit' title='Edit URL' onClick={editUrl(index, url)}/>
                            <Icon link fitted name='delete' title='Remove URL'  onClick={removeUrl(index)}/>
                          </Label.Detail>
                        }
                      </Label>
                    );
                  })}
                  {canEditDetails && <Label as='a' color='green' onClick={createUrl}><Icon name='plus' />Add URL</Label>}
                </Label.Group>
              </Form.Field>
            </Form.Group>
          </Form>
          <Box icon='settings' title='Details' ref='details' stacked={Boolean(projectId)}>
            <Form error={Boolean(errors.details.length)}>
              <Form.Group widths='two'>
                <Form.Input readOnly={!canEditDetails} label='Domain' value={project.domain || ''} onChange={this.handleChange}
                    type='text' name='domain' autoComplete='off' placeholder='Project Domain' width='eight' error={errors.fields['domain']}
                />
                {this.renderDropdown('projectManager', 'Project Manager', project.projectManager, 'Select Project Manager...', 'eight', users, isEntitiesFetching, errors, !canEditDetails)}
              </Form.Group>

              <Form.Group widths='two'>
                {this.renderDropdown('serviceCenter', 'Service Center', project.serviceCenter, 'Select Service Center...', 'eight', serviceCenters, isEntitiesFetching, errors, !canEditDetails)}
                {this.renderDropdown('businessUnit', 'Business Unit', project.businessUnit, 'Select Business Unit...', 'eight', businessUnits, isEntitiesFetching, errors, !canEditDetails)}
              </Form.Group>
              <Message error list={errors.details}/>
            </Form>
          </Box>
          <Divider hidden/>
          {this.renderServices(project, services, fetching, !canEditMatrix)}
          {canEditMatrix && <Button color='green' icon='save' title='Save project' labelPosition='left' content='Save Project' onClick={this.handleSubmit} className='floating' size='big' />}
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
  onCreateUrl: React.PropTypes.func.isRequired,
  onEditUrl: React.PropTypes.func.isRequired,
  onRemoveUrl: React.PropTypes.func.isRequired,
  onSave: React.PropTypes.func.isRequired,
  onDelete: React.PropTypes.func.isRequired,
  canEditDetails: React.PropTypes.bool,
};

const mapStateToProps = (state, ownProps) => {
  const auth = state.auth;
  const paramId = ownProps.params.id;
  const projects = state.projects;
  const project = projects.selected;
  const emptyProject = { matrix: [], urls: [] };
  const isFetching = paramId && (paramId !== project.id || project.isFetching);
  const isServicesFetching = state.services.isFetching;
  const users = Object.values(state.users.items);
  const authUser = auth.user;
  const selectedProject = { ...emptyProject, ...projects.items[paramId] };
  selectedProject.urls = selectedProject.urls || [];

  let entities = Object.values(state.entities.items);

  // The only entities we show for a RI and a CP are:
  // * the entities assigned to the RI
  // * the businessUnit and serviceCenter assigned to the current project
  if (authUser.role !== AUTH_ADMIN_ROLE) {
    entities = entities.filter(entity =>
      authUser.entities
        .concat(selectedProject.businessUnit)
        .concat(selectedProject.serviceCenter)
        .includes(entity.id));
  }
  
  const userEntities = authUser.entities || [];
  const commonEntities = userEntities.find(id => [selectedProject.businessUnit, selectedProject.serviceCenter].includes(id)) || [];
  // Details of the project can be edited if the user is an admin
  // or if the user is a RI and it's a project linked to that user
  // or if the user is a RI and it's a new project
  const canEditDetails = authUser.role === AUTH_ADMIN_ROLE || (authUser.role === AUTH_RI_ROLE && (commonEntities.length > 0 || !project.id));
  
  const services = groupByPackage(state.services.items);
  return {
    auth,
    project: selectedProject,
    isFetching,
    projectId: paramId,
    businessUnits: getEntitiesAsOptions({ entities: getByType(entities, 'businessUnit') }),
    serviceCenters: getEntitiesAsOptions({ entities: getByType(entities, 'serviceCenter') }),
    users: getUsersAsOptions(users),
    isEntitiesFetching: state.entities.isFetching,
    services,
    isServicesFetching,
    canEditDetails
  };
};

const mapDispatchToProps = dispatch => ({
  fetchProject: id => dispatch(ProjectsThunks.fetch(id)),
  fetchEntities: () => dispatch(EntitiesThunks.fetchIfNeeded()),
  fetchServices: () => dispatch(ServicesThunks.fetchIfNeeded()),
  fetchUsers: () => dispatch(UsersThunks.fetchIfNeeded()),
  onCreateUrl: (id) => {
    const cb = (url) => dispatch(ProjectsActions.addUrl(id, url));
    dispatch(ModalActions.openNewUrlModal(cb));
  },
  onEditUrl: (id, index, url) => {
    const cb = (url) => dispatch(ProjectsActions.editUrl(id, index, url));
    dispatch(ModalActions.openEditUrlModal(url, cb));
  },
  onRemoveUrl: (id, index) => dispatch(ProjectsActions.removeUrl(id, index)),
  onSave: project => dispatch(ProjectsThunks.save(project, ToastsActions.savedProjectSuccessNotification(project.name))),
  onDelete: project => {
    const del = () => dispatch(ProjectsThunks.delete(project, push('/projects')));
    dispatch(ModalActions.openRemoveProjectModal(project, del));
  }
});

const ProjectPage = connect(
  mapStateToProps,
  mapDispatchToProps
)(ProjectComponent);

export default ProjectPage;
