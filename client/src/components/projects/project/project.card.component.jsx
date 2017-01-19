// React
import React from 'react';
import { Link } from 'react-router';
import { Card, Icon, Label } from 'semantic-ui-react';


import './project.card.component.scss';

// ProjectCard Component
class ProjectCard extends React.Component {

  render = () => {
    const { project } = this.props;
    return (
      <Card className='project-card' raised>
        <Card.Content>
          <Card.Header as='h4'title={project.name} className='ui left floated link'>
            <Link to={`/projects/${project.id}`}><Icon fitted name='travel' />{project.name.toUpperCase()}</Link>
          </Card.Header>
          <Label as='a' href={project.url} color='blue' className='ui right floated'>
            <Icon name='linkify' />
            URL
          </Label>
        </Card.Content>
        <Card.Content extra >
          <div className='domain' >
            {project.domain}
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
