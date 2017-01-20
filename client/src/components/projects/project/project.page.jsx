// React
import React from 'react';
import { Link } from 'react-router';
import { connect } from 'react-redux';
import { push } from 'react-router-redux';
import { Button, Container, Divider, Form, Header, Icon, Table, Segment } from 'semantic-ui-react';
import classnames from 'classnames';

import Matrix from './matrix/matrix.component';
import Box from '../../common/box.component';

// Thunks / Actions
import ProjectsThunks from '../../../modules/projects/projects.thunks';
import OrganizationsThunks from '../../../modules/organizations/organizations.thunks';
import ServicesThunks from '../../../modules/services/services.thunks';

import { getOrganizationsAsOptions, getByType } from '../../../modules/organizations/organizations.selectors';
import { groupByPackage } from '../../../modules/services/services.selectors';

// Style
import './project.page.scss';

// Project Component
class ProjectComponent extends React.Component {

  state = { project: {} }

  componentWillMount = () => {
    this.setState({ project: { ...this.props.project } });
  }

  componentWillReceiveProps = (nextProps) => {
    this.setState({ project: { ...nextProps.project } });
  }

  componentDidMount = () => {
    const { projectId } = this.props;
    Promise.all([
      this.props.fetchOrganizations(),
      this.props.fetchServices()
    ]).then(()=>{
      this.props.fetchProject(projectId);
    });
  }

  handleChange = (e, { name, value, checked }) => {
    const { project } = this.state;
    const state = {
      project: { ...project, [name]:value || checked }
    };
    this.setState(state);
  }

  handleSubmit = (e) => {
    e.preventDefault();
    const stateProject = this.state.project;
    const project = { ...stateProject };
    this.props.onSave(project);
  }

  renderServices = (project, services) => {
    const matrix = {};
    project.matrix.forEach((m) => matrix[m.service] = m);
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
              return <Matrix key={service.id} matrix={matrix[service.id]} service={service} />;
            })}
          </Table.Body>
        </Table>
      );
    });
  }

  render = () => {
    const { isFetching, serviceCenters, entities, isOrganizationsFetching, services } = this.props;
    const { project } = this.state;
    const readOnly = false;
    const classes = classnames({ readonly: readOnly });
    return (
      <Container className='project-page'>
        <Segment loading={isFetching} padded>
          <Header as='h1'>
            <Link to={'/projects'}>
              <Icon name='arrow left' fitted/>
            </Link>
            {project.name}
            {project.url && <Button as='a' href={project.url} content='URL' icon='linkify' labelPosition='left' color='blue' floated='right' />}
          </Header>
          <Divider hidden/>
          <Box icon='settings' title='Details' stacked>
            <Form>
              <Form.Group>
                <Form.Input readOnly={readOnly} label='Name' value={project.name || ''} onChange={this.handleChange}
                  type='text' name='name' autoComplete='off' placeholder='Project name' width='four'
                />
                <Form.Input readOnly={readOnly} label='Domain' value={project.domain || ''} onChange={this.handleChange}
                    type='text' name='domain' autoComplete='off' placeholder='Project domain' width='four'
                />
                <Form.Input readOnly={readOnly} label='URL' value={project.url || ''} onChange={this.handleChange}
                    type='text' name='url' autoComplete='off' placeholder='Docktor group url'  width='eight'
                />
              </Form.Group>

              <Form.Group widths='two'>
                <Form.Dropdown disabled={readOnly} placeholder='Select service center...' fluid search selection loading={isOrganizationsFetching}
                  label='Service Center' name='serviceCenter' options={serviceCenters} value={project.serviceCenter || ''} onChange={this.handleChange} className={classes}
                />
                <Form.Dropdown disabled={readOnly} placeholder='Select entity...' fluid search selection loading={isOrganizationsFetching}
                  label='Entity' name='entity' options={entities} value={project.entity || ''} onChange={this.handleChange} className={classes}
                />
              </Form.Group>
            </Form>
          </Box>
          <Divider hidden/>
          {this.renderServices(project, services)}
          <Divider hidden/>
          <Button fluid color='green' content='Save' loading={isFetching} onClick={this.handleSubmit} />
        </Segment>
      </Container>
    );
  }
}

ProjectComponent.propTypes = {
  project: React.PropTypes.object,
  isFetching: React.PropTypes.bool,
  entities: React.PropTypes.array,
  serviceCenters: React.PropTypes.array,
  isOrganizationsFetching: React.PropTypes.bool,
  services: React.PropTypes.object,
  projectId: React.PropTypes.string.isRequired,
  fetchProject: React.PropTypes.func.isRequired,
  fetchOrganizations: React.PropTypes.func.isRequired,
  fetchServices: React.PropTypes.func.isRequired,
  onSave: React.PropTypes.func.isRequired
};

const mapStateToProps = (state, ownProps) => {
  const paramId = ownProps.params.id;
  const projects = state.projects;
  const project = projects.selected;
  const emptyProject = { matrix: [] };
  const isFetching = paramId && (paramId !== project.id || project.isFetching);
  const organizations = Object.values(state.organizations.items);
  const services = groupByPackage(state.services.items);
  return {
    project: { ...emptyProject, ...projects.items[project.id] },
    isFetching,
    projectId: paramId,
    entities: getOrganizationsAsOptions(getByType(organizations, 'entity')),
    serviceCenters: getOrganizationsAsOptions(getByType(organizations, 'serviceCenter')),
    isOrganizationsFetching: state.organizations.isFetching,
    services
  };
};

const mapDispatchToProps = dispatch => ({
  fetchProject: id => dispatch(ProjectsThunks.fetch(id)),
  fetchOrganizations: () => dispatch(OrganizationsThunks.fetchIfNeeded()),
  fetchServices: () => dispatch(ServicesThunks.fetchIfNeeded()),
  onSave: project => dispatch(ProjectsThunks.save(project, push('/projects')))
});

const ProjectPage = connect(
  mapStateToProps,
  mapDispatchToProps
)(ProjectComponent);

export default ProjectPage;
