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
    const { isFetching } = this.props;
    const { project } = this.state;
    return (
      <Container className='project-page'>
        <Segment loading={isFetching} padded>
          <Header as='h1'>
            <Link to={'/projects'}>
              <Icon name='arrow left' fitted/>
            </Link>
            {project.name}
          </Header>
          <Divider hidden/>
          <Form onSubmit={this.handleSubmit}>
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
  projectId: React.PropTypes.string.isRequired,
  fetchProject: React.PropTypes.func.isRequired,
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
