// React
import React from 'react';
import PropTypes from 'prop-types';
import { Link } from 'react-router';
import { connect } from 'react-redux';
import { push } from 'react-router-redux';
import DocumentTitle from 'react-document-title';
import { Button, Container, Divider, Form, Grid, Icon, Label, List, Message, Table, Segment, Popup } from 'semantic-ui-react';
import Joi from 'joi-browser';

import Matrix from './matrix/matrix.component';
import Box from '../../common/box.component';

// Thunks / Actions
import ProjectsThunks from '../../../modules/projects/projects.thunks';
import EntitiesThunks from '../../../modules/entities/entities.thunks';
import { fetchIndicators } from '../../../modules/indicators/indicators.thunks';
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

import { AUTH_DEPUTY_ROLE, AUTH_PM_ROLE, AUTH_RI_ROLE, AUTH_ADMIN_ROLE } from '../../../modules/auth/auth.constants';

// Style
import './project.page.scss';

const Legend = ({ options, status }) => (
  <Box icon='help circle' title='Color Legend'>
    <Divider horizontal>Progress & Goal</Divider>
    <Grid columns={2} relaxed>
      <Grid.Column>
        {/*Next line is used to separate options list in two parts, we use Math.ceil to make the left side bigger than the right one*/}
        {options.slice(0, Math.ceil(options.length / 2)).map((opt) => (
          <List.Item key={opt.value}>
            <Label color={opt.label.color} horizontal>{opt.text}</Label>
            {opt.title}
          </List.Item>
        ))}
      </Grid.Column>
      <Grid.Column>
        {options.slice(Math.ceil(options.length / 2)).map((opt) => (
          <List.Item key={opt.value}>
            <Label color={opt.label.color} horizontal>{opt.text}</Label>
            {opt.title}
          </List.Item>
        ))}
      </Grid.Column>
    </Grid>

    <Divider horizontal>Indicator</Divider>
    <Grid columns={2} relaxed>
      <Grid.Column>
        {status.slice(0, Math.ceil(status.length / 2)).map((stat) => (
          <List.Item key={stat.value}>
            <Label className='status-label' circular empty color={stat.color} />
            <span>{stat.title}</span>
          </List.Item>
        ))}
      </Grid.Column>
      <Grid.Column>
        {status.slice(Math.ceil(status.length / 2)).map((stat) => (
          <List.Item key={stat.value}>
            <Label className='status-label' circular empty color={stat.color} />
            <span>{stat.title}</span>
          </List.Item>
        ))}
      </Grid.Column>
    </Grid>
  </Box>
);
Legend.propTypes = {
  options: PropTypes.array,
  status: PropTypes.array
};

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
      { text: 'Git (GitLab)', value: 'GitLab' },
      { text: 'Git (Other)', value: 'Git' },
      { text: 'Mercurial', value: 'Mercurial' },
      { text: 'CVS', value: 'CVS' },
      { text: 'TFS', value: 'TFS' },
      { text: 'Other', value: 'Other' }
    ],
    technologies: []
  }

  schema = Joi.object().keys({
    name: Joi.string().trim().required().label('Project Name'),
    client: Joi.string().trim().empty('').label('Client'),
    docktorGroupURL: Joi.string().trim().empty('').label('Docktor URL'),
    mode: Joi.string().trim().empty('').label('Mode'),
    deliverables: Joi.boolean().label('Deliverables'),
    isCDKApplicable: Joi.boolean().label('The CDK is not applicable globally'),
    sourceCode: Joi.boolean().label('Source Code'),
    specifications: Joi.boolean().label('Specifications'),
    projectManager: Joi.string().trim().alphanum().empty('').label('Project Manager'),
    serviceCenter: Joi.string().trim().alphanum().empty('').label('Service Center'),
    businessUnit: Joi.string().trim().alphanum().empty('').label('Business Unit'),
    explanation: Joi.string().trim().empty('').label('Explanation'),
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
        this.props.fetchIndicators(projectId);
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

  renderPackages = (packages, indicators, isFetching, isConnectedUserAdmin, readonly, isIsolatedNetwork) => {
    if (isFetching) {
      return <p>Fetching Matrix...</p>;
    }

    const packageList = Object.entries(packages).map(([pckg, servicesList]) => (
      <Table key={pckg} celled striped compact>
        <Table.Header>
          <Table.Row>
            <Table.HeaderCell width='7'>{pckg}</Table.HeaderCell>
            <Table.HeaderCell width='1'>Deployed</Table.HeaderCell>
            <Table.HeaderCell width='2'>Progress</Table.HeaderCell>
            <Table.HeaderCell width='2'>Goal</Table.HeaderCell>
            <Table.HeaderCell width='1'>Priority</Table.HeaderCell>
            <Table.HeaderCell width='2'>Due Date</Table.HeaderCell>
            <Table.HeaderCell width='1'>Comment</Table.HeaderCell>
          </Table.Row>
        </Table.Header>
        <Table.Body>
          {servicesList.map((service) => (
            <Matrix
              key={service.id} readOnly={readonly} isConnectedUserAdmin={isConnectedUserAdmin}
              serviceId={service.id} matrix={this.state.matrix[service.id] || {}} service={service}
              indicators={indicators} onChange={this.handleMatrix} isIsolatedNetwork={isIsolatedNetwork}
            />
          ))}
        </Table.Body>
      </Table>
    ));

    const title = this.props.projectId ? 'D.A.D - Project ' + (this.props.project && this.props.project.name) : 'D.A.D - New Project';
    return <DocumentTitle title={title}><div>{packageList}</div></DocumentTitle>;
  }

  renderDropdown = (name, label, value, placeholder, options, isFetching, errors, canEdit) => {
    if (!canEdit) {
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

  renderConsolidationCriteriaField = (domain, canEdit, errors) => {
    const selectedCriterias = domain || [];
    const options = selectedCriterias.map((d) => ({ text: d, value: d }));
    if (!canEdit) {
      return (
        <div className='field'>
          <label>Consolidation criteria</label>
          <div>
            {selectedCriterias.map((criteria) => <Label key={criteria} size='large'>{criteria}</Label>)}
          </div>
        </div>
      );
    }
    return (
      <Form.Dropdown
        label='Consolidation criteria' placeholder='Rennes, Offshore, ...' fluid multiple selection allowAdditions
        onChange={this.handleChange}
        name='domain' search value={selectedCriterias} options={options} error={errors && errors.fields['domain']}
      />
    );
  }

  renderMultipleSearchSelectionDropdown = (name, label, selectedValuesDropDown = [], values, placeholder, canEdit) => {
    selectedValuesDropDown = selectedValuesDropDown || [];
    if (!canEdit) {
      return (
        <Form.Field>
          <label>{label}</label>
          {values
            .filter((element) => selectedValuesDropDown.includes(element.value))
            .map((element) => <Label key={element.value} size='large'>{element.text}</Label>)
          }
        </Form.Field>
      );
    }
    return (
      <Form.Dropdown
        label={label} placeholder={placeholder} fluid multiple selection onChange={this.handleChange}
        name={name} allowAdditions={true} search value={selectedValuesDropDown} options={values}
      />
    );
  }

  render = () => {
    const {
      isFetching, serviceCenters, businessUnits, indicators,
      isEntitiesFetching, services, isServicesFetching,
      users, projectId, isAdmin, isRI, isPM, isDeputy
    } = this.props;
    const { project, errors } = this.state;
    const fetching = isFetching || isServicesFetching;

    // The list of technologies options must contain the default technologies *and* the custom technologies
    // added by the user. If we don't do the concatenation, the component won't be able to display the
    // user-defined technologies. We dedupe the technologies using a Set.
    const technologiesOptions =
      Array
        .from(new Set(this.props.technologies.concat(project.technologies || [])))
        .map((technology) => ({ text: technology, value: technology }));

    // Remove the 'None' user from the list of users because the list of deputies doesn't have a default option with empty value
    const usersWithoutNone = users.filter((user) => user && user.text !== 'None');
    return (
      <Container className='project-page'>
        <Segment loading={fetching} padded>
          <Form>
            <h1 className='layout horizontal center justified'>
              <Link to={'/projects'}>
                <Icon name='arrow left' fitted />
              </Link>

              <Form.Input
                className='flex projectName' value={project.name || ''} onChange={this.handleChange} type='text' name='name'
                placeholder='Project Name' error={errors.fields['name']} readOnly={isPM || isDeputy}
              />

              {/*Only admins and RIs can delete a project*/}
              {((isAdmin || isRI) && typeof projectId !== 'undefined') && <Button color='red' icon='trash' labelPosition='left' title='Delete project' content='Delete Project' onClick={this.handleRemove} />}
            </h1>

            <Divider hidden />

            <Form.Group>
              <Form.TextArea
                label='Description' value={project.description || ''} onChange={this.handleChange} autoHeight type='text'
                name='description' autoComplete='off' placeholder='Project description' width='sixteen' error={errors.fields['description']}
              />
            </Form.Group>
          </Form>

          <Box icon='settings' title='Details' ref='details' stacked={Boolean(projectId)}>
            <Form error={Boolean(errors.details.length)}>
              <Grid columns='equal' divided>
                <Grid.Row>
                  <Grid.Column>
                    <h3>Project Data</h3>

                    {this.renderDropdown('projectManager', 'Project Manager', project.projectManager, 'Select Project Manager...', users, isEntitiesFetching, errors, (isAdmin || isRI))}

                    {this.renderMultipleSearchSelectionDropdown('deputies', 'Deputies', project.deputies || [], usersWithoutNone, 'Add deputy...', (isAdmin || isRI))}

                    <Form.Input
                      label='Client' value={project.client || ''} onChange={this.handleChange} type='text' name='client'
                      autoComplete='on' placeholder='Project Client' error={errors.fields['client']} readOnly={isPM || isDeputy}
                    />

                    {/*The field Domain was renamed Consolidation criteria only in the GUI. All references named Domain in code is corresponding to the Consolidation criteria field*/}
                    <Popup trigger={
                      this.renderConsolidationCriteriaField(project.domain, (isAdmin || isRI), errors)
                    } position='top right' wide size='mini' on='click' inverted>
                      <Popup.Content>
                        Useful to add your own filters (for the Search Options and in the Export).
                        Several values allowed, press Enter to validate.
                      </Popup.Content>
                    </Popup>

                    {this.renderDropdown('serviceCenter', 'Service Center', project.serviceCenter, 'Select Service Center...', serviceCenters, isEntitiesFetching, errors, (isAdmin || isRI))}

                    {this.renderDropdown('businessUnit', 'Business Unit', project.businessUnit, 'Select Business Unit...', businessUnits, isEntitiesFetching, errors, (isAdmin || isRI))}

                    {/*Only admins are allowed to set the Docktor URL*/}
                    <Form.Input
                      label='Docktor Group URL' value={project.docktorGroupURL || ''} onChange={this.handleChange} type='text' name='docktorGroupURL' autoComplete='on'
                      placeholder='http://<DocktorURL>/#!/groups/<GroupId>' error={errors.fields['docktorGroupURL']} readOnly={!isAdmin}
                    />
                  </Grid.Column>

                  <Grid.Column>
                    <h3>Technical Data</h3>

                    {this.renderMultipleSearchSelectionDropdown('technologies', 'Technologies', project.technologies || [], technologiesOptions, 'Java, .NET...', (isAdmin || isRI))}

                    {/*The deployment mode is editable by the admins only*/}
                    {this.renderDropdown('mode', 'Deployment Mode', project.mode, 'SaaS, DMZ...', this.state.modes, false, errors, isAdmin)}

                    <h4>Version Control</h4>

                    <Form.Checkbox readOnly={isPM || isDeputy} label='Deliverables' name='deliverables' checked={Boolean(project.deliverables)} onChange={this.handleChange} />

                    <Form.Checkbox readOnly={isPM || isDeputy} label='Source Code' name='sourceCode' checked={Boolean(project.sourceCode)} onChange={this.handleChange} />

                    <Form.Checkbox readOnly={isPM || isDeputy} label='Specifications' name='specifications' checked={Boolean(project.specifications)} onChange={this.handleChange} />

                    {this.renderDropdown('versionControlSystem', 'Version Control System', project.versionControlSystem, 'SVN, Git...', this.state.versionControlSystems, false, errors, (isAdmin || isRI))}
                    <div className='ui divider' />
                    <h4>Applicability of CDK</h4>
                    <div className='ui segment' title='WARNING: The entire matrix will be disabled'>
                      <Form.Checkbox readOnly={isPM || isDeputy} label='The CDK is not applicable globally' name='isCDKApplicable' checked={Boolean(project.isCDKApplicable)} disabled={typeof projectId !== 'undefined'} onChange={this.handleChange} />
                    </div>
                    {project.isCDKApplicable && <Form.TextArea readOnly={isPM || isDeputy} label='Explanation' name='explanation' value={project.explanation || ''} placeholder='The CDK is not applicable because...' disabled={typeof projectId !== 'undefined'} onChange={this.handleChange} />}
                  </Grid.Column>
                </Grid.Row>
              </Grid>

              <Message error list={errors.details} />
            </Form>
          </Box>

          <Legend options={options} status={status} />

          <Divider hidden />

          {this.renderPackages(services, indicators, fetching, isAdmin, this.state.project.isCDKApplicable, this.state.project.mode === 'Isolated Network')}

          <Button
            color='green' icon='save' title='Save project' labelPosition='left' content='Save Project'
            onClick={this.handleSubmit} className='floating' size='big'
          />
        </Segment>
      </Container >
    );
  }
}

ProjectComponent.propTypes = {
  auth: PropTypes.object.isRequired,
  project: PropTypes.object,
  isFetching: PropTypes.bool,
  businessUnits: PropTypes.array,
  serviceCenters: PropTypes.array,
  explanation: PropTypes.array,
  indicators: PropTypes.object,
  isEntitiesFetching: PropTypes.bool,
  users: PropTypes.array,
  services: PropTypes.object,
  technologies: PropTypes.array,
  isServicesFetching: PropTypes.bool,
  projectId: PropTypes.string,
  fetchProject: PropTypes.func.isRequired,
  fetchEntities: PropTypes.func.isRequired,
  fetchIndicators: PropTypes.func.isRequired,
  fetchServices: PropTypes.func.isRequired,
  fetchUsers: PropTypes.func.isRequired,
  fetchTechnologies: PropTypes.func.isRequired,
  onSave: PropTypes.func.isRequired,
  onDelete: PropTypes.func.isRequired,
  isAdmin: PropTypes.bool,
  isRI: PropTypes.bool,
  isCDKApplicable: PropTypes.bool,
  isPM: PropTypes.bool,
  isDeputy: PropTypes.bool
};

const mapStateToProps = (state, ownProps) => {
  const auth = state.auth;
  const paramId = ownProps.params.id;
  const projects = state.projects;
  const project = projects.selected;
  const indicators = state.indicators;
  const technologies = state.technologies.items;
  const emptyProject = { matrix: [] };
  const isFetching = paramId && (paramId !== project.id || project.isFetching);
  const isServicesFetching = state.services.isFetching;
  const users = Object.values(state.users.items);
  const authUser = auth.user;
  const selectedProject = { ...emptyProject, ...projects.items[paramId] };

  let entities = Object.values(state.entities.items);

  // The only entities we show for a RI, a PM and a Deputy are:
  // * the entities assigned to the RI
  // * the businessUnit and serviceCenter assigned to the current project
  if (authUser.role !== AUTH_ADMIN_ROLE) {
    entities = entities.filter((entity) =>
      authUser.entities
        .concat(selectedProject.businessUnit || [])
        .concat(selectedProject.serviceCenter || [])
        .includes(entity.id));
  }
  const services = groupByPackage(state.services.items);

  const isAdmin = authUser.role === AUTH_ADMIN_ROLE;
  const isRI = authUser.role === AUTH_RI_ROLE;
  const isPM = authUser.role === AUTH_PM_ROLE;
  const isDeputy = authUser.role === AUTH_DEPUTY_ROLE;
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
    indicators,
    isServicesFetching,
    isAdmin,
    isRI,
    isPM,
    isDeputy
  };
};

const mapDispatchToProps = (dispatch) => ({
  fetchProject: (id) => dispatch(ProjectsThunks.fetch(id)),
  fetchEntities: () => dispatch(EntitiesThunks.fetchIfNeeded()),
  fetchIndicators: (id) => dispatch(fetchIndicators(id)),
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
