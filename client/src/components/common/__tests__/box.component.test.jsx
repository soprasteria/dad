import React from 'react';
import { shallow } from 'enzyme';

import Box from '../box.component';

describe('<Box />', () => {
  it('should render a .settings icon', () => {
    const wrapper = shallow(<Box icon='settings' />);
    expect(wrapper.find({ name: 'settings' })).toHaveLength(1);
  });
});
