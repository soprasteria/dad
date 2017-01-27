// React
import React from 'react';
import { Link } from 'react-router';
import { Card, Icon, Label } from 'semantic-ui-react';


import './project.card.component.scss';

// ProjectCard Component
class ProjectCard extends React.Component {

  render = () => {
    const { project } = this.props;
    project.matrix = project.matrix || [];
    const filteredMatrix = project.matrix.filter(m => m.goal != -1);
    const goals = filteredMatrix.map(m => [m.progress, m.goal])
      .reduce((acc, [progress, goal]) => {
        if (goal === -1) {return 0;}
        return progress >= goal ? acc + 1 : acc;
      }, 0);
    return (
      <Card className='project-card' raised>
        <Card.Content>
          <Card.Header as='h4'title={project.name} className='ui left floated link'>
            <Link to={`/projects/${project.id}`}><Icon fitted name='travel' />{project.name.toUpperCase()}</Link>
          </Card.Header>
          <Label color='blue' image title={`${goals} goal(s) reached`} className='ui right floated'>
            <Icon fitted name='star' />
            <Label.Detail>
              {`${goals}/${filteredMatrix.length}`}
            </Label.Detail>
          </Label>
        </Card.Content>
        <Card.Content extra >
          <div className='domain' >
            {project.domain || 'No Domain'}
          </div>
        </Card.Content>
      </Card>
    );
  }
}

ProjectCard.propTypes = {
  project: React.PropTypes.object
};

export default ProjectCard;
