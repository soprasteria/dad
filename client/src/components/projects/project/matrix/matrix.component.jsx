// React
import React from 'react';
import DebounceInput from 'react-debounce-input';
import { Form, Table, Button, Icon, Popup, Label } from 'semantic-ui-react';
import ReactDatePicker from 'react-datepicker';
import classNames from 'classnames';
import moment from 'moment';

import 'react-datepicker/dist/react-datepicker.css';

import { options, priorities, status } from '../../../../modules/services/services.constants';

import './matrix.component.scss';

// Matrix Component Component
class Matrix extends React.Component {

  handleChange = (e, { name, value }) => {
    this.props.onChange(this.props.serviceId, { ...this.props.matrix, [name]: value });
  }

  handleChangeComment = ({ target }) => {
    this.props.onChange(this.props.serviceId, { ...this.props.matrix, comment: target.value });
  }

  handleChangeDueDate = (date) => {
    this.props.onChange(this.props.serviceId, { ...this.props.matrix, dueDate: date || undefined });
  }

  render = () => {
    const { service, matrix, indicators, readOnly } = this.props;
    return (
      <Table.Row className='matrix-component'>
        {this.renderCells(service, matrix, indicators.items, readOnly)}
      </Table.Row>
    );
  }

  renderCells = (service, matrix, indicators, readOnly) => {
    matrix.progress = typeof matrix.progress === 'number' ? matrix.progress : -1;
    matrix.goal = typeof matrix.goal === 'number' ? matrix.goal : -1;
    matrix.priority = typeof matrix.priority === 'string' && matrix.priority !== '' ? matrix.priority : 'N/A';

    const serviceStatus = this.getServiceStatus(service, indicators);
    const progressOption = options.find((elm) => elm.value === matrix.progress);
    const priorityOption = priorities.find((elm) => elm.value === matrix.priority);
    const goalOption = options.find((elm) => elm.value === matrix.goal);
    const dueDate = matrix.dueDate ? moment(matrix.dueDate) : '';
    const expandComment = this.state && this.state.expandComment;
    const serviceNameCell = (
      <Table.Cell key='service'>
        {/* If serviceStatus is in an unknown status, the label indicator will not be visible for users */}
        <Label
          className={classNames({ invisible: !serviceStatus }, 'status-label')} circular
          title={serviceStatus ? serviceStatus.title : ''}
          color={serviceStatus ? serviceStatus.color : 'grey'}
          />
        <span>
          {service.name}
        </span>
      </Table.Cell>
    );

    const setExpandComment = (expandComment) => this.setState((prevState) => {
      return {
        ...prevState,
        expandComment
      };
    });

    if (expandComment) {
      // When the comment is expanded, the only 2 cells are the service name and the comment
      return [
        serviceNameCell,
        (<Table.Cell key='comment' colSpan={5}>
          <Form>
            <DebounceInput autoFocus readOnly={readOnly} debounceTimeout={600} element={Form.TextArea} autoHeight
              placeholder={readOnly ? '' : 'Add a comment'} name='comment' value={matrix.comment}
              onChange={this.handleChangeComment} onBlur={() => setExpandComment(false)}
              />
          </Form>
        </Table.Cell>)
      ];
    } else {
      return [
        serviceNameCell,
        (<Table.Cell key='progress'>
          <Form>
            {readOnly
              ? (<div>{progressOption.text}</div>)
              : (<Form.Dropdown placeholder='Progress' fluid selection name='progress' title={progressOption.title}
                options={options} value={matrix.progress} onChange={this.handleChange} className={progressOption.label.color}
                />)
            }
          </Form>
        </Table.Cell>),
        (<Table.Cell key='goal'>
          <Form>
            {readOnly
              ? (<div>{goalOption.text}</div>)
              : (<Form.Dropdown placeholder='Goal' fluid selection name='goal' title={goalOption.title}
                options={options} value={matrix.goal} onChange={this.handleChange} className={goalOption.label.color}
                />)
            }
          </Form>
        </Table.Cell>),
        (<Table.Cell key='priority'>
          <Form>
            {readOnly
              ? (<div>{priorityOption.text}</div>)
              : (<Form.Dropdown placeholder='Priority' fluid selection name='priority' title={priorityOption.title}
                options={priorities} value={matrix.priority} onChange={this.handleChange}
                />)
            }
          </Form>
        </Table.Cell>),
        (<Table.Cell key='dueDate'>
          <Form>
            <ReactDatePicker dateFormat='DD/MM/YYYY' placeholderText='DD/MM/YYYY' selected={dueDate} onChange={this.handleChangeDueDate} />
          </Form>
        </Table.Cell>),
        (<Table.Cell key='comment' className={classNames(readOnly, 'comment', 'center')}>
          <Form>
            <Popup
              trigger={
                <Button icon name='comment' onClick={() => setExpandComment(true)} color={matrix.comment ? 'blue' : null}>
                  <Icon name='comment' />
                </Button>
              }
              content={matrix.comment ? matrix.comment : 'Click to add a comment'}
              header={matrix.comment ? 'Click to edit' : null}
              inverted
              />
          </Form>
        </Table.Cell>)
      ];
    }
  }

  // Method used to get the status indicator for a functional service. It can determine which service has the best status to display only this one
  getServiceStatus = (service, indicators) => {
    var serviceStatus = undefined;
    // Condition used to test if the functionnal service has some technical services associated
    if (service.services) {
      const indicatorsTable = Object.values(indicators);
      var matchingIndicators = service.services.map((technicalServiceName) =>
        indicatorsTable.filter((indic) => indic.service === technicalServiceName)
      );
      // Converting tables of tables to tables of objects 
      matchingIndicators = [].concat(...matchingIndicators);
      // compare 2 indicators status and return the best one
      serviceStatus = matchingIndicators.reduce((currentStatus, indicator) => {
        const newStatus = status.find((s) => s.text === (indicator && indicator.status));
        if (currentStatus && newStatus && currentStatus.value > newStatus.value) {
          return currentStatus;
        } else if (currentStatus && !newStatus) {
          return currentStatus;
        } else {
          return newStatus;
        }
      }, undefined);
    }
    return serviceStatus;
  }
}

Matrix.propTypes = {
  serviceId: React.PropTypes.string,
  indicators: React.PropTypes.object,
  matrix: React.PropTypes.object,
  service: React.PropTypes.object,
  onChange: React.PropTypes.func,
  readOnly: React.PropTypes.bool
};

export default Matrix;
