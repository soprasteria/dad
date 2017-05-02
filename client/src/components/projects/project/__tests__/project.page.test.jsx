import React from 'react';
import { shallow } from 'enzyme';

import { ProjectComponent } from '../project.page';
import { AUTH_ADMIN_ROLE } from '../../../../modules/auth/auth.constants';

function setup(customProps = {}) {
  const props = {
    auth: {
      user: {
        role: AUTH_ADMIN_ROLE,
        entities: []
      }
    },
    canEditDetails: true,
    technologies: [],
    users: [],
    serviceCenters: [],
    businessUnits: [],
    services: {},
    project: {
      matrix: []
    },
    fetchProject: jest.fn(),
    fetchEntities: jest.fn(),
    fetchServices: jest.fn(),
    fetchUsers: jest.fn(),
    fetchTechnologies: jest.fn(),
    onSave: jest.fn(),
    onDelete: jest.fn(),
    ...customProps,
  };

  const wrapper = shallow(<ProjectComponent {...props} />);

  return {
    props,
    wrapper
  };
}

describe('<ProjectComponent />', () => {
  it('renders two <h3 /> components', () => {
    const { wrapper } = setup();
    expect(wrapper.find('h3').length).toEqual(2);
  });

  describe('with an empty matrix', () => {
    describe('with canEditDetails set to true', () => {
      it('should display x <Form.Dropdown />', () => {
        const { wrapper } = setup();
        expect(wrapper.find('FormDropdown').length).toEqual(6);
      });
    });

    describe('with canEditDetails set to false', () => {
      it('should display no <Form.Dropdown />', () => {
        const { wrapper } = setup({ canEditDetails: false });
        expect(wrapper.find('FormDropdown').length).toEqual(0);
      });
    });
  });
});
