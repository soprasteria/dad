// React
import React from 'react';

import { Form, Table } from 'semantic-ui-react';

import { options } from '../../../../modules/services/services.constants';

import './matrix.component.scss';

// Matrix Component Component
class Matrix extends React.Component {

  state = { matrix : {} }

  componentWillMount = () => {
    this.setState({ matrix : { ...this.props.matrix } });
  }

  componentWillReceiveProps = (nextProps) => {
    this.setState({ matrix : { ...nextProps.matrix } });
  }

  handleChange = (e, { name, value }) => {
    this .setState({ matrix: { ...this.state.matrix, [name]:value } });
  }

  render = () => {
    const { service } = this.props;
    const { matrix } = this.state;
    return (
      <Table.Row>
        <Table.Cell>{service.name}</Table.Cell>
        <Table.Cell>
          <Form>
            <Form.Dropdown placeholder='Progress' fluid selection name='progress'
              options={options} value={matrix.progress || -1} onChange={this.handleChange}
            />
          </Form>
        </Table.Cell>
        <Table.Cell>
          <Form>
            <Form.Dropdown placeholder='Goal' fluid selection name='goal'
              options={options} value={matrix.goal || -1} onChange={this.handleChange}
            />
          </Form>
        </Table.Cell>
        <Table.Cell>
          <Form>
            <Form.TextArea placeholder='Add a comment' name='comment' autoHeight value={matrix.comment} onChange={this.handleChange} />
          </Form>
        </Table.Cell>
      </Table.Row>
    );
  }
}

Matrix.propTypes = {
  matrix: React.PropTypes.object,
  service: React.PropTypes.object
};

export default Matrix;
