// React
import React from 'react';

import { Dropdown, List } from 'semantic-ui-react';

import './matrix.component.scss';

// Matrix Component Component
class Matrix extends React.Component {

  render = () => {
    const { service } = this.props;
    return (
      <List.Item>
        <List.Content>
          <List.Header>{service.name}</List.Header>
          <Dropdown defaultValue='Progress' placeholder='Select progress'/>
          <Dropdown defaultValue='Goal' placeholder='Select goal'/>
        </List.Content>
      </List.Item>
    );
  }
}

Matrix.propTypes = {
  matrix: React.PropTypes.object,
  service: React.PropTypes.object
};

export default Matrix;
