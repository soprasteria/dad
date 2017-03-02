// React
import React from 'react';
import DebounceInput from 'react-debounce-input';
import { Form, Table } from 'semantic-ui-react';

import { options, priorities } from '../../../../modules/services/services.constants';

import './matrix.component.scss';

// Golang date format
const dueDateSuffix = 'T00:00:00Z';
const defaultDueDate = '0001-01-01' + dueDateSuffix;

// Matrix Component Component
class Matrix extends React.Component {

  handleChange = (e, { name, value }) => {
    this.props.onChange(this.props.serviceId, { ...this.props.matrix, [name]:value });
  }

  handleChangeComment = ({ target }) => {
    this.props.onChange(this.props.serviceId, { ...this.props.matrix, comment:target.value });
  }

  handleChangeDueDate = ({ target }) => {
    this.props.onChange(this.props.serviceId, { ...this.props.matrix, dueDate: target.value ? target.value + dueDateSuffix : '' });
  }

  render = () => {
    const { service, matrix, readOnly } = this.props;
    matrix.progress = typeof matrix.progress === 'number' ? matrix.progress : -1;
    matrix.goal = typeof matrix.goal === 'number' ? matrix.goal : -1;
    matrix.priority = typeof matrix.priority === 'string' && matrix.priority !== '' ? matrix.priority : 'N/A';

    const progressOption = options.find(elm => elm.value === matrix.progress);
    const priorityOption = priorities.find(elm => elm.value === matrix.priority);
    const goalOption = options.find(elm => elm.value === matrix.goal);
    const dueDate = matrix.dueDate && matrix.dueDate !== defaultDueDate ? matrix.dueDate.substr(0, 10) : '';
    return (
      <Table.Row className='matrix-component'>
        <Table.Cell>{service.name}</Table.Cell>
        <Table.Cell>
          <Form>
            {readOnly
              ? (<div>{progressOption.text}</div>)
              : (<Form.Dropdown placeholder='Progress' fluid selection name='progress' title={progressOption.title}
                  options={options} value={matrix.progress} onChange={this.handleChange} className={progressOption.label.color}
                />)
            }
          </Form>
        </Table.Cell>
        <Table.Cell>
          <Form>
            {readOnly
              ? (<div>{goalOption.text}</div>)
              : (<Form.Dropdown placeholder='Goal' fluid selection name='goal' title={goalOption.title}
                  options={options} value={matrix.goal} onChange={this.handleChange} className={goalOption.label.color}
                />)
            }
          </Form>
        </Table.Cell>
        <Table.Cell>
          <Form>
            {readOnly
              ? (<div>{priorityOption.text}</div>)
              : (<Form.Dropdown placeholder='Priority' fluid selection name='priority' title={priorityOption.title}
                options={priorities} value={matrix.priority} onChange={this.handleChange}
              />)
            }
          </Form>
        </Table.Cell>
        <Table.Cell>
          <Form>
            <input type='date' name='dueDate' value={dueDate} onChange={this.handleChangeDueDate}/>
          </Form>
        </Table.Cell>
        <Table.Cell>
          <Form>
            <DebounceInput readOnly={readOnly} debounceTimeout={600} element={Form.TextArea} autoHeight
              placeholder={readOnly ? '' : 'Add a comment'} name='comment' value={matrix.comment}
              onChange={this.handleChangeComment}
            />
          </Form>
        </Table.Cell>
      </Table.Row>
    );
  }
}

Matrix.propTypes = {
  serviceId: React.PropTypes.string,
  matrix: React.PropTypes.object,
  service: React.PropTypes.object,
  onChange: React.PropTypes.func,
  readOnly: React.PropTypes.bool
};

export default Matrix;
