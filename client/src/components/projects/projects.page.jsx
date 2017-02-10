// React
import React from 'react';
import { connect } from 'react-redux';
import { Link } from 'react-router';
import { Button, Card, Container, Icon, Input, Label, Segment } from 'semantic-ui-react';
import DebounceInput from 'react-debounce-input';
import { AUTH_CP_ROLE } from '../../modules/auth/auth.constants';

// API Fetching
import ProjectsThunks from '../../modules/projects/projects.thunks';
import EntitiesThunks from '../../modules/entities/entities.thunks';
import ProjectsActions from '../../modules/projects/projects.actions';

// Selectors
import { getFilteredProjects } from '../../modules/projects/projects.selectors';

// Components
import ProjectCard from './project/project.card.component';

import './projects.page.scss';

//Site Component using react-leaflet
class Projects extends React.Component {

  componentWillMount = () => {
    Promise.all([this.props.fetchEntities()]).then(() => {
      this.props.fetchProjects();
    });
  }

  renderCards = (projects, entities) => {
    if (projects.length) {
      return (
        <Card.Group className='centered'>
          {projects.map(project => {
            return (
              <ProjectCard project={project} key={project.id} businessUnit={entities[project.businessUnit] || {}} serviceCenter={entities[project.serviceCenter] || {}}/>
            );
          })}
        </Card.Group>
      );
    }
    return <p>No projects found.</p>;
  }

  render = () => {
    const { projects, entities, filterValue, isFetching, changeFilter, auth } = this.props;
    return (
      <Container fluid className='projects-page'>
        <Segment.Group raised>
          <Segment clearing>
            <Input icon labelPosition='left corner'>
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
            {auth.user.role !== AUTH_CP_ROLE && <Button as={Link} content='New Project' icon='plus' labelPosition='left' color='green' floated='right' to={'/projects/new'} />}
          </Segment>
          <Segment loading={isFetching}>
            {this.renderCards(projects, entities)}
          </Segment>
        </Segment.Group>
      </Container>
    );
  }
}

Projects.propTypes = {
  auth: React.PropTypes.object.isRequired,
  projects: React.PropTypes.array,
  entities: React.PropTypes.object,
  filterValue: React.PropTypes.string,
  isFetching: React.PropTypes.bool,
  fetchProjects: React.PropTypes.func.isRequired,
  fetchEntities: React.PropTypes.func.isRequired,
  changeFilter: React.PropTypes.func.isRequired
};

// Function to map state to container props
const mapStateToProjectsProps = (state) => {
  const filterValue = state.projects.filterValue;
  const entities = state.entities.items;
  const projects = getFilteredProjects(state.projects.items, entities, filterValue);
  const isFetching = state.projects.isFetclhing || state.entities.isFetching;
  return {
    filterValue,
    projects,
    entities,
    isFetching,
    auth: state.auth
  };
};

// Function to map dispatch to container props
const mapDispatchToProjectsProps = (dispatch) => {
  return {
    fetchProjects : () => dispatch(ProjectsThunks.fetchIfNeeded()),
    fetchEntities : () => dispatch(EntitiesThunks.fetchIfNeeded()),
    changeFilter: filterValue => dispatch(ProjectsActions.changeFilter(filterValue))
  };
};

// Redux container to Sites component
const ProjectsPage = connect(
  mapStateToProjectsProps,
  mapDispatchToProjectsProps
)(Projects);

export default ProjectsPage;
