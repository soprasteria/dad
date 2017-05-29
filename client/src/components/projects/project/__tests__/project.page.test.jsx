import React from 'react';
import { shallow } from 'enzyme';

import { ProjectComponent } from '../project.page';
import { AUTH_ADMIN_ROLE } from '../../../../modules/auth/auth.constants';

const defaultProps = {
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
};

function setup(customProps = {}) {
  const props = {
    ...defaultProps,
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

describe('renderTechnologiesField', () => {
  it('renders a drop down with technologies', () => {
    const projectComponent = new ProjectComponent(defaultProps);
    const result = projectComponent.renderTechnologiesField(['java', '.NET', 'Pega'], ['java', '.NET'], false);
    expect(result).not.toBeNull();
  });
  it('renders div with technologies in readonly mode', () => {
    const projectComponent = new ProjectComponent(defaultProps);
    const result = projectComponent.renderTechnologiesField(['java', '.NET', 'Pega'], ['java', '.NET'], true);
    expect(result).not.toBeNull();
  });
  it('renders a drop down with empty technologies', () => {
    const projectComponent = new ProjectComponent(defaultProps);
    const result = projectComponent.renderTechnologiesField([], ['java', '.NET'], false);
    expect(result).not.toBeNull();
  });
  it('renders div with empty technologies in readonly mode', () => {
    const projectComponent = new ProjectComponent(defaultProps);
    const result = projectComponent.renderTechnologiesField([], ['java', '.NET'], true);
    expect(result).not.toBeNull();
  });
  it('renders a drop down with null technologies', () => {
    const projectComponent = new ProjectComponent(defaultProps);
    const result = projectComponent.renderTechnologiesField(null, ['java', '.NET'], false);
    expect(result).not.toBeNull();
  });
  it('renders div with null technologies in readonly mode', () => {
    const projectComponent = new ProjectComponent(defaultProps);
    const result = projectComponent.renderTechnologiesField(null, ['java', '.NET'], true);
    expect(result).not.toBeNull();
  });
  it('renders a drop down with undefined technologies', () => {
    const projectComponent = new ProjectComponent(defaultProps);
    const result = projectComponent.renderTechnologiesField(undefined, ['java', '.NET'], false);
    expect(result).not.toBeNull();
  });
  it('renders div with undefined technologies in readonly mode', () => {
    const projectComponent = new ProjectComponent(defaultProps);
    const result = projectComponent.renderTechnologiesField(undefined, ['java', '.NET'], true);
    expect(result).not.toBeNull();
  });
});

describe('renderConsolidationCriteriaField', () => {
  it('renders a drop down with consolidation criterias', () => {
    const projectComponent = new ProjectComponent(defaultProps);
    const result = projectComponent.renderConsolidationCriteriaField(['Rennes', 'BIOS'], false, null);
    expect(result).not.toBeNull();
  });
  it('renders div with consolidation criterias in readonly mode', () => {
    const projectComponent = new ProjectComponent(defaultProps);
    const result = projectComponent.renderTechnologiesField(['Rennes', 'BIOS'], true, null);
    expect(result).not.toBeNull();
  });
  it('renders a drop down with empty consolidation criterias', () => {
    const projectComponent = new ProjectComponent(defaultProps);
    const result = projectComponent.renderTechnologiesField([], false, null);
    expect(result).not.toBeNull();
  });
  it('renders div with empty consolidation criterias in readonly mode', () => {
    const projectComponent = new ProjectComponent(defaultProps);
    const result = projectComponent.renderTechnologiesField([], true, null);
    expect(result).not.toBeNull();
  });
  it('renders a drop down with null consolidation criterias', () => {
    const projectComponent = new ProjectComponent(defaultProps);
    const result = projectComponent.renderTechnologiesField(null, false, null);
    expect(result).not.toBeNull();
  });
  it('renders div with null consolidation criterias in readonly mode', () => {
    const projectComponent = new ProjectComponent(defaultProps);
    const result = projectComponent.renderTechnologiesField(null, true, null);
    expect(result).not.toBeNull();
  });
  it('renders a drop down with undefined consolidation criterias', () => {
    const projectComponent = new ProjectComponent(defaultProps);
    const result = projectComponent.renderTechnologiesField(undefined, false, null);
    expect(result).not.toBeNull();
  });
  it('renders div with undefined consolidation criterias in readonly mode', () => {
    const projectComponent = new ProjectComponent(defaultProps);
    const result = projectComponent.renderTechnologiesField(undefined, true, null);
    expect(result).not.toBeNull();
  });
});
