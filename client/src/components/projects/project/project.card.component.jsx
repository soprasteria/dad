// React
import React from 'react';
import { Link } from 'react-router';
import { Card } from 'semantic-ui-react';


import './project.card.component.scss';

// ProjectCard Component
class ProjectCard extends React.Component {

  render = () => {
    const { project } = this.props;
    return (
      <Card className='project-card' raised>
        <Card.Content>
          <Link to={`/projects/${project.id}`}>
            {project.name}
          </Link>
        </Card.Content>
        <Card.Content extra>
          <div className='domain'>
            <p>{project.domain}</p>
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
