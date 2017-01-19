// React
import React from 'react';
import { Link } from 'react-router';
import { connect } from 'react-redux';
import { push } from 'react-router-redux';
import { Button, Container, Divider, Form, Header, Icon, Segment } from 'semantic-ui-react';

// Thunks / Actions
import ProjectsThunks from '../../../modules/projects/projects.thunks';

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
    this.props.fetchProject(projectId);
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

  render = () => {
    const { isFetching, serviceCenters, entities, isOrganizationsFetching } = this.props;
    const { project } = this.state;
    return (
      <Container className='project-page'>
        <Segment loading={isFetching} padded>
          <Header as='h1'>
            <Link to={'/projects'}>
              <Icon name='arrow left' fitted/>
            </Link>
            {project.name}
            <Button content='Docktor URL' icon='linkify' labelPosition='left' color='blue' floated='right' />
          </Header>
          <Divider hidden/>
          <Form onSubmit={this.handleSubmit}>
            <Form.Group widths='two' >
              <Form.Input readOnly={false} label='Name' value={project.name || ''} onChange={this.handleChange}
                type='text' name='name' autoComplete='off' placeholder='Project name'
              />
              <Form.Input readOnly={false} label='Domain' value={project.domain || ''} onChange={this.handleChange}
                  type='text' name='domain' autoComplete='off' placeholder='Project domain'
              />
            </Form.Group>

            <Form.Input readOnly={false} label='url' value={project.url || ''} onChange={this.handleChange}
                type='text' name='url' autoComplete='off' placeholder='Docktor group url'
            />

            <Form.Group widths='two'>
              <Form.Dropdown readOnly={false} placeholder='Select entity...' fluid search selection loading={isOrganizationsFetching}
                name='entity' options={entities} value={project.entity || []} onChange={this.handleChange}
              />
              <Form.Dropdown readOnly={false} placeholder='Select service center...' fluid search selection loading={isOrganizationsFetching}
                name='serviceCenter' options={serviceCenters} value={project.serviceCenter || []} onChange={this.handleChange}
              />
            </Form.Group>
            <Button fluid color='green' content='Save' loading={isFetching} />
          </Form>
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
  projectId: React.PropTypes.string.isRequired,
  fetchProject: React.PropTypes.func.isRequired,
  fetchOrganizations: React.PropTypes.func.isRequired,
  onSave: React.PropTypes.func.isRequired
};

const mapStateToProps = (state, ownProps) => {
  const paramId = ownProps.params.id;
  const projects = state.projects;
  const project = projects.selected;
  const emptyProject = {};
  const isFetching = paramId && (paramId !== project.id || project.isFetching);
  return {
    project: projects.items[project.id] || emptyProject,
    isFetching,
    projectId: paramId
  };
};

const mapDispatchToProps = dispatch => ({
  fetchProject: id => dispatch(ProjectsThunks.fetch(id)),
  onSave: project => dispatch(ProjectsThunks.save(project, push('/projects')))
});

const ProjectPage = connect(
  mapStateToProps,
  mapDispatchToProps
)(ProjectComponent);

export default ProjectPage;
