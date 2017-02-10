// React
import React from 'react';
import { Link } from 'react-router';
import { Card, Icon, Label } from 'semantic-ui-react';


import './project.card.component.scss';

// ProjectCard Component
class ProjectCard extends React.Component {

  render = () => {
    const { project, businessUnit, serviceCenter } = this.props;
    project.matrix = project.matrix || [];
    const filteredMatrix = project.matrix.filter(m => m.goal !== -1);
    const goals = filteredMatrix.map(m => [m.progress, m.goal])
      .reduce((acc, [progress, goal]) => {
        if (progress === -1) {progress = 0;}
        const res = acc  + Math.min(progress * 100 / goal, 100);
        return res;
      }, 0);
    return (
      <Card className='project-card' raised>
        <Card.Content>
          <Card.Header as='h4'title={project.name} className='ui left floated link'>
            <Link to={`/projects/${project.id}`}><Icon fitted name='travel' />{project.name.toUpperCase()}</Link>
          </Card.Header>
          <Label color='blue' image title={'Goal completion rates'} className='ui right floated'>
            <Icon fitted name='line chart' />
            <Label.Detail>
              {Math.floor(goals / filteredMatrix.length) + '%'}
            </Label.Detail>
          </Label>
          <Card.Meta className='domain' title={project.domain}>
            {project.domain || 'No Domain'}
          </Card.Meta>
        </Card.Content>
        <Card.Content extra >
          <div className='left floated service-center' title={serviceCenter.name}>
            {serviceCenter.name || 'No Service Center'}
          </div>
          <div className='right floated business-unit' title={businessUnit.name}>
            {businessUnit.name || 'No Business Unit'}
          </div>
        </Card.Content>
      </Card>
    );
  }
}

ProjectCard.propTypes = {
  project: React.PropTypes.object,
  businessUnit: React.PropTypes.object,
  serviceCenter: React.PropTypes.object,
};

export default ProjectCard;
