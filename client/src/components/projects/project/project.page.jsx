// React
import React from 'react';
import { Link } from 'react-router';
import { connect } from 'react-redux';
import { push } from 'react-router-redux';
import DocumentTitle from 'react-document-title';
import { Button, Container, Divider, Form, Grid, Icon, Label, List, Message, Table, Segment } from 'semantic-ui-react';

import Joi from 'joi-browser';

import Matrix from './matrix/matrix.component';
import Box from '../../common/box.component';

// Thunks / Actions
import ProjectsThunks from '../../../modules/projects/projects.thunks';
import EntitiesThunks from '../../../modules/entities/entities.thunks';
import TechnologiesThunks from '../../../modules/technologies/technologies.thunks';
import ServicesThunks from '../../../modules/services/services.thunks';
import { options, status } from '../../../modules/services/services.constants';
import UsersThunks from '../../../modules/users/users.thunks';
import ModalActions from '../../../modules/modal/modal.actions';
import ToastsActions from '../../../modules/toasts/toasts.actions';

import { getEntitiesAsOptions, getByType } from '../../../modules/entities/entities.selectors';
import { groupByPackage } from '../../../modules/services/services.selectors';
import { flattenTechnologies } from '../../../modules/technologies/technologies.selectors';
import { getUsersAsOptions } from '../../../modules/users/users.selectors';

import { parseError } from '../../../modules/utils/forms';

import { AUTH_CP_ROLE, AUTH_RI_ROLE, AUTH_ADMIN_ROLE } from '../../../modules/auth/auth.constants';

// Style
import './project.page.scss';

// Project Component
export class ProjectComponent extends React.Component {

  state = {
    errors: {
      details: [],
      fields: {}
    },
    project: {},
    matrix: {},
    modes: [
      { text: '', value: '' },
      { text: 'SaaS', value: 'SaaS' },
      { text: 'DMZ', value: 'DMZ' },
      { text: 'Isolated Network', value: 'Isolated Network' }
    ],
    versionControlSystems: [
      { text: '', value: '' },
      { text: 'SVN', value: 'SVN' },
      { text: 'Git', value: 'Git' },
      { text: 'Mercurial', value: 'Mercurial' },
      { text: 'CVS', value: 'CVS' },
      { text: 'TFS', value: 'TFS' },
      { text: 'Other', value: 'Other' }
    ],
    technologies: []
  }

  schema = Joi.object().keys({
    name: Joi.string().trim().required().label('Project Name'),
    domain: Joi.string().trim().empty('').label('Domain'),
    client: Joi.string().trim().empty('').label('Client'),
    mode: Joi.string().trim().empty('').label('Mode'),
    deliverables: Joi.boolean().label('Deliverables'),
    sourceCode: Joi.boolean().label('Source Code'),
    specifications: Joi.boolean().label('Specifications'),
    projectManager: Joi.string().trim().alphanum().empty('').label('Project Manager'),
    serviceCenter: Joi.string().trim().alphanum().empty('').label('Service Center'),
    businessUnit: Joi.string().trim().alphanum().empty('').label('Business Unit')
  }).or('serviceCenter', 'businessUnit').label('Service Center or Business Unit');

  componentWillMount = () => {
    const matrix = {};
    const project = this.props.project;
    project.matrix.forEach((m) => matrix[m.service] = m);
    this.setState({ project: { ...project }, errors: { details: [], fields: {} }, matrix });
  }

  componentWillReceiveProps = (nextProps) => {
    const project = nextProps.project;
    if (!project.isEditing) {
      const matrix = {};
      project.matrix.forEach((m) => matrix[m.service] = m);
      this.setState({ project: { ...project }, errors: { details: [], fields: {} }, matrix });
    } else {
      this.setState({ project: { ...this.state.project } });
    }
  }

  componentDidMount = () => {
    const { projectId } = this.props;
    Promise.all([
      this.props.fetchEntities(),
      this.props.fetchServices(),
      this.props.fetchUsers(),
      this.props.fetchTechnologies()
    ]).then(() => {
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

  handleChange = (e, { name, value, checked }) => {
    const { project, errors } = this.state;
    const state = {
      project: {
        ...project,
        // "checked" is used for checkboxes, because their "value" doesn't change
        [name]: value || checked
      },
      errors: {
        details: [...errors.details],
        fields: {
          ...errors.fields
        }
      }
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
    const result = Joi.validate(this.state.project, this.schema, { abortEarly: false, allowUnknown: true });
    if (result.error) {
      const errors = parseError(result.error);
      if (errors.fields['Service Center or Business Unit']) {
        errors.fields.serviceCenter = true;
        errors.fields.businessUnit = true;
        delete errors.fields['Service Center or Business Unit'];
      }
      window.scrollTo(0, 0);
      this.refs.details.setState({ stacked: false });
      this.setState({ errors: errors });
    }
    return result;
  }

  handleSubmit = (e) => {
    e.preventDefault();
    const formValidationResult = this.isFormValid();
    if (!Boolean(formValidationResult.error)) {
      const { matrix } = this.state;
      // Use the "project" object returned by Joi.validate instead of `this.state.project` because
      // Joi converts the type of the checkbox values to boolean automatically
      const project = formValidationResult.value;
      const modifiedProject = {
        ...project,
        matrix: Object.values(matrix)
      };
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

    const packageList = Object.entries(services).map(([pckg, servicesList]) => {
      return (
        <Table key={pckg} celled striped compact>
          <Table.Header>
            <Table.Row>
              <Table.HeaderCell width='8'>{pckg}</Table.HeaderCell>
              <Table.HeaderCell width='2'>Progress</Table.HeaderCell>
              <Table.HeaderCell width='2'>Goal</Table.HeaderCell>
              <Table.HeaderCell width='1'>Priority</Table.HeaderCell>
              <Table.HeaderCell width='2'>Due Date</Table.HeaderCell>
              <Table.HeaderCell width='1'>Comment</Table.HeaderCell>
            </Table.Row>
          </Table.Header>
          <Table.Body>
            {servicesList.map((service) => {
              return <Matrix readOnly={readOnly} serviceId={service.id} key={service.id} matrix={this.state.matrix[service.id] || {}} service={service} onChange={this.handleMatrix} />;
            })}
          </Table.Body>
        </Table>
      );
    });

    const title = this.props.projectId ? 'D.A.D - Project ' + (this.props.project && this.props.project.name) : 'D.A.D - New Project';
    return <DocumentTitle title={title}><div>{packageList}</div></DocumentTitle>;
  }

  renderDropdown = (name, label, value, placeholder, options, isFetching, errors, readOnly) => {
    if (readOnly) {
      const option = options.find((elm) => elm.value === value);
      return (
        <Form.Input readOnly label={label} value={(option && option.text) || ''} onChange={this.handleChange}
          type='text' autoComplete='off' placeholder={`No ${label}`}
        />
      );
    }
    return (
      <Form.Dropdown placeholder={placeholder} fluid search selection loading={isFetching}
        label={label} name={name} options={options} value={value || ''} onChange={this.handleChange} error={errors.fields[name]}
      />
    );
  }

  renderTechnologiesField = (selectedTechnologies = [], technologies, readOnly) => {
    selectedTechnologies = selectedTechnologies || [];
    if (readOnly) {
      return (
        <div>
          {selectedTechnologies.map((technology) => <Label size='large'>{technology}</Label>)}
        </div>
      );
    }
    return (
      <Form.Dropdown
        label='Technologies' placeholder='Java, .NET...' fluid multiple selection onChange={this.handleChange}
        name='technologies' allowAdditions={true} search value={selectedTechnologies} options={technologies}
      />
    );
  }

  render = () => {
    const {
      isFetching, serviceCenters, businessUnits,
      isEntitiesFetching, services, isServicesFetching,
      users, projectId, canEditDetails
    } = this.props;

    const { project, errors } = this.state;
    const fetching = isFetching || isServicesFetching;
    const authUser = this.props.auth.user;
    const canEditMatrix = canEditDetails || (authUser.role === AUTH_CP_ROLE && project.projectManager === authUser.id);

    // The list of technologies options must contain the default technologies *and* the custom technologies
    // added by the user. If we don't do the concatenation, the component won't be able to display the
    // user-defined technologies. We dedupe the technologies using a Set.
    const technologiesOptions =
      Array
        .from(new Set(this.props.technologies.concat(project.technologies || [])))
        .map((technology) => ({ text: technology, value: technology }));

    return (
      <Container className='project-page'>

        <Segment loading={fetching} padded>

          <Form>
            <h1 className='layout horizontal center justified'>
              <Link to={'/projects'}>
                <Icon name='arrow left' fitted />
              </Link>
              <Form.Input className='flex projectName' readOnly={!canEditDetails} value={project.name || ''} onChange={this.handleChange}
                type='text' name='name' autoComplete='off' placeholder='Project Name' error={errors.fields['name']}
              />
              {(!isFetching && canEditDetails && projectId !== null) && <Button color='red' icon='trash' labelPosition='left' title='Delete project' content='Delete Project' onClick={this.handleRemove} />}
            </h1>

            <Divider hidden />
            <Form.Group>
              <Form.TextArea
                readOnly={!canEditMatrix} label='Description' value={project.description || ''} onChange={this.handleChange} autoHeight
                type='text' name='description' autoComplete='off' placeholder='Project description' width='sixteen' error={errors.fields['description']}
              />
            </Form.Group>
          </Form>

          <Box icon='settings' title='Details' ref='details' stacked={Boolean(projectId)}>
            <Form error={Boolean(errors.details.length)}>
              <Grid columns='equal' divided>
                <Grid.Row>
                  <Grid.Column>
                    <h3>Project Data</h3>
                    {this.renderDropdown('projectManager', 'Project Manager', project.projectManager, 'Select Project Manager...', users, isEntitiesFetching, errors, !canEditDetails)}
                    <Form.Input readOnly={!canEditDetails} label='Client' value={project.client || ''} onChange={this.handleChange}
                      type='text' name='client' autoComplete='on' placeholder='Project Client' error={errors.fields['client']}
                    />
                    <Form.Input readOnly={!canEditDetails} label='Domain' value={project.domain || ''} onChange={this.handleChange}
                      type='text' name='domain' autoComplete='on' placeholder='Project Domain' error={errors.fields['domain']}
                    />
                    {this.renderDropdown('serviceCenter', 'Service Center', project.serviceCenter, 'Select Service Center...', serviceCenters, isEntitiesFetching, errors, !canEditDetails)}
                    {this.renderDropdown('businessUnit', 'Business Unit', project.businessUnit, 'Select Business Unit...', businessUnits, isEntitiesFetching, errors, !canEditDetails)}
                  </Grid.Column>

                  <Grid.Column>
                    <h3>Technical Data</h3>

                    {this.renderTechnologiesField(project.technologies || [], technologiesOptions, !canEditDetails)}

                    {this.renderDropdown('mode', 'Deployment Mode', project.mode, 'SaaS, DMZ...', this.state.modes, false, errors, !canEditDetails)}

                    <h4>Version Control</h4>
                    <Form.Checkbox readOnly={!canEditDetails} label='Deliverables' name='deliverables'
                      checked={Boolean(project.deliverables)} onChange={this.handleChange} />

                    <Form.Checkbox readOnly={!canEditDetails} label='Source Code' name='sourceCode'
                      checked={Boolean(project.sourceCode)} onChange={this.handleChange} />

                    <Form.Checkbox readOnly={!canEditDetails} label='Specifications' name='specifications'
                      checked={Boolean(project.specifications)} onChange={this.handleChange} />

                    {this.renderDropdown('versionControlSystem', 'Version Control System', project.versionControlSystem, 'SVN, Git...', this.state.versionControlSystems, false, errors, !canEditDetails)}
                  </Grid.Column>
                </Grid.Row>
              </Grid>

              <Message error list={errors.details} />
            </Form>
          </Box>
          <Box icon='help circle' title='Color Legend' ref='legend'>
            <Divider horizontal> Maturity Legend </Divider>
            <Grid columns={2} relaxed>
              <Grid.Column>
                {options.slice(0, Math.ceil(options.length / 2)).map((opt) => {
                  return (
                    <List.Item key={opt.value}>
                      <Label color={opt.label.color} horizontal>{opt.text}</Label>
                      {opt.title}
                    </List.Item>
                  );
                })}
              </Grid.Column>
              <Grid.Column>
                {options.slice(Math.ceil(options.length / 2)).map((opt) => {
                  return (
                    <List.Item key={opt.value}>
                      <Label color={opt.label.color} horizontal>{opt.text}</Label>
                      {opt.title}
                    </List.Item>
                  );
                })}
              </Grid.Column>
            </Grid>
            <Divider horizontal> Indicator Legend </Divider>
            <Grid columns={2} relaxed>
              <Grid.Column>
                {status.slice(0, Math.ceil(status.length / 2)).map((stat) => {
                  return (
                    <List.Item key={stat.value}>
                      <Label className='status-label' circular empty color={stat.color}/>
                      <span>{stat.title}</span>
                    </List.Item>
                  );
                })}
              </Grid.Column>
              <Grid.Column>
                {status.slice(Math.ceil(status.length / 2)).map((stat) => {
                  return (
                      <List.Item key={stat.value}>
                        <Label className='status-label' circular empty color={stat.color}/>
                        <span>{stat.title}</span>
                      </List.Item>
                  );
                })}
              </Grid.Column>
            </Grid>
          </Box>
          <Divider hidden />
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
  technologies: React.PropTypes.array,
  isServicesFetching: React.PropTypes.bool,
  projectId: React.PropTypes.string,
  fetchProject: React.PropTypes.func.isRequired,
  fetchEntities: React.PropTypes.func.isRequired,
  fetchServices: React.PropTypes.func.isRequired,
  fetchUsers: React.PropTypes.func.isRequired,
  fetchTechnologies: React.PropTypes.func.isRequired,
  onSave: React.PropTypes.func.isRequired,
  onDelete: React.PropTypes.func.isRequired,
  canEditDetails: React.PropTypes.bool,
};

const mapStateToProps = (state, ownProps) => {
  const auth = state.auth;
  const paramId = ownProps.params.id;
  const projects = state.projects;
  const project = projects.selected;
  const technologies = state.technologies.items;
  const emptyProject = { matrix: [] };
  const isFetching = paramId && (paramId !== project.id || project.isFetching);
  const isServicesFetching = state.services.isFetching;
  const users = Object.values(state.users.items);
  const authUser = auth.user;
  const selectedProject = { ...emptyProject, ...projects.items[paramId] };

  let entities = Object.values(state.entities.items);
  const userEntities = authUser.entities || [];

  // The only entities we show for a RI and a CP are:
  // * the entities assigned to the RI
  // * the businessUnit and serviceCenter assigned to the current project
  if (authUser.role !== AUTH_ADMIN_ROLE) {
    entities = entities.filter((entity) =>
      authUser.entities
        .concat(selectedProject.businessUnit || [])
        .concat(selectedProject.serviceCenter || [])
        .includes(entity.id));
  }

  const commonEntities = userEntities.find((id) => [selectedProject.businessUnit, selectedProject.serviceCenter].includes(id)) || [];

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
    technologies: flattenTechnologies(technologies),
    isEntitiesFetching: state.entities.isFetching,
    services,
    isServicesFetching,
    canEditDetails
  };
};

const mapDispatchToProps = (dispatch) => ({
  fetchProject: (id) => dispatch(ProjectsThunks.fetch(id)),
  fetchEntities: () => dispatch(EntitiesThunks.fetchIfNeeded()),
  fetchServices: () => dispatch(ServicesThunks.fetchIfNeeded()),
  fetchUsers: () => dispatch(UsersThunks.fetchIfNeeded()),
  fetchTechnologies: () => dispatch(TechnologiesThunks.fetchIfNeeded()),
  onSave: (project) => dispatch(ProjectsThunks.save(project, (id) => push('/projects/' + id), ToastsActions.savedSuccessNotification('Project ' + project.name))),
  onDelete: (project) => {
    const del = () => dispatch(ProjectsThunks.delete(project, push('/projects')));
    dispatch(ModalActions.openRemoveProjectModal(project, del));
  }
});

const ProjectPage = connect(
  mapStateToProps,
  mapDispatchToProps
)(ProjectComponent);

export default ProjectPage;
