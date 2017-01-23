// React
import React from 'react';
import { connect } from 'react-redux';
import { Link } from 'react-router';
import { Button, Card, Container, Icon, Input, Label, Segment } from 'semantic-ui-react';
import DebounceInput from 'react-debounce-input';

// API Fetching
import ProjectsThunks from '../../modules/projects/projects.thunks';
import ProjectsActions from '../../modules/projects/projects.actions';

// Selectors
import { getFilteredProjects } from '../../modules/projects/projects.selectors';

// Components
import ProjectCard from './project/project.card.component';

import './projects.page.scss';

//Site Component using react-leaflet
class Projects extends React.Component {

  componentWillMount = () => {
    this.props.fetchProjects();
  }

  renderCards = (projects) => {
    if (projects.length) {
      return (
        <Card.Group>
          {projects.map(project => {
            return (
              <ProjectCard project={project} key={project.id} />
            );
          })}
        </Card.Group>
      );
    }
    return <p>No projects found.</p>;
  }

  render = () => {
    const { projects, filterValue, isFetching, changeFilter } = this.props;
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
                <Button as={Link} content='New Project' icon='plus' labelPosition='left' color='green' floated='right' to={'/projects/new'} />
          </Segment>
          <Segment loading={isFetching}>
            {this.renderCards(projects)}
          </Segment>
        </Segment.Group>
      </Container>
    );
  }
}

Projects.propTypes = {
  projects: React.PropTypes.array,
  filterValue: React.PropTypes.string,
  isFetching: React.PropTypes.bool,
  fetchProjects: React.PropTypes.func.isRequired,
  changeFilter: React.PropTypes.func.isRequired
};

// Function to map state to container props
const mapStateToProjectsProps = (state) => {
  const filterValue = state.projects.filterValue;
  const projects = getFilteredProjects(state.projects.items, filterValue);
  const isFetching = state.projects.isFetching;
  return { filterValue, projects, isFetching };
};

// Function to map dispatch to container props
const mapDispatchToProjectsProps = (dispatch) => {
  return {
    fetchProjects : () => dispatch(ProjectsThunks.fetchIfNeeded()),
    changeFilter: filterValue => dispatch(ProjectsActions.changeFilter(filterValue))
  };
};

// Redux container to Sites component
const ProjectsPage = connect(
  mapStateToProjectsProps,
  mapDispatchToProjectsProps
)(Projects);

export default ProjectsPage;
