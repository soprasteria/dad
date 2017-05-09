// React
import React from 'react';
import { Form, Label } from 'semantic-ui-react';

import './editable.link.scss';

// EditableLink
export class EditableLink extends React.Component {

  render = () => {
    const { labelValue, title, url } = this.props;
    return (
    <Form.Field>
      <label>{labelValue}</label>
      {this.renderEditableLink(title, url)}
    </Form.Field>
    );
  }

  renderEditableLink = (title, url) => {
    if (title && url) {
      return (
        <div>
          <Label className='label-ui' id='underline-url' size='large'>
            {<a className='truncated-url' href={url} title={title} target='_blank'>{title}</a>}
          </Label>
        </div>
      );
    }
    if (url) {
      return (
        <div>
          <Label className='label-ui' id='underline-url' size='large'>
            {<a className='truncated-url' href={url} title={title} target='_blank'>{url}</a>}
          </Label>
        </div>
      );
    }
    return (
        <div>
          <Label className='label-ui' size='large'>+ Add Docktor URL</Label>
        </div>
    );
  }
}

EditableLink.propTypes = {
  labelValue: React.PropTypes.string.isRequired,
  title: React.PropTypes.string,
  url: React.PropTypes.string
};

export default EditableLink;
