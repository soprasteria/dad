// React
import React from 'react';
import { Link } from 'react-router';
import { Card, Icon, Label } from 'semantic-ui-react';
import classNames from 'classnames';

import { calculateProgress } from '../../../modules/utils/projects';
import './project.card.component.scss';

// ProjectCard Component
class ProjectCard extends React.Component {

  render = () => {
    const { project, businessUnit, serviceCenter } = this.props;
    project.matrix = project.matrix || [];
    const filteredMatrixGoals = project.matrix.filter((m) => m.goal >= 0);
    const filteredMatrixProgress = project.matrix.filter((m) => m.progress >= 0);
    const goalMessage = (filteredMatrixGoals.length === 0 && filteredMatrixProgress.length === 0) ? '-' :  // If nothing was specified, then the output is ''-'
          filteredMatrixGoals.length === 0 ? 'N/A' : Math.floor(calculateProgress(project)) + '%'; // Else if there is no goal specified the output is N/A, otherwise we do the maths.
    const domainClassnames = classNames({ filled: project.domain }, 'domain');
    const serviceCenterClassnames = classNames({ filled: serviceCenter.name }, 'left floated service-center');
    const businessUnitClassnames = classNames({ filled: businessUnit.name }, 'right floated business-unit');
    return (
      <Card className='project-card' raised>
        <Card.Content>
          <Card.Header as='h4'title={project.name} className='ui left floated link'>
            <Link to={`/projects/${project.id}`}><Icon fitted name='travel' />{project.name.toUpperCase()}</Link>
          </Card.Header>
          <Label color='blue' image title={'Goal completion rates'} className='ui right floated'>
            <Icon fitted name='line chart' />
            <Label.Detail>
              {goalMessage}
            </Label.Detail>
          </Label>
          <Card.Meta className={domainClassnames} title={project.domain}>
            {project.domain || 'No Domain'}
          </Card.Meta>
        </Card.Content>
        <Card.Content extra >
          <div className={serviceCenterClassnames} title={serviceCenter.name}>
            {serviceCenter.name || 'No Service Center'}
          </div>
          <div className={businessUnitClassnames} title={businessUnit.name}>
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
