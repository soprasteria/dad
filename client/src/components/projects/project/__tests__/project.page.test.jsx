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
  fetchIndicators: jest.fn(),
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
    describe('with isAdmin set to true', () => {
      it('should display x <Form.Dropdown />', () => {
        const { wrapper } = setup({ isAdmin: true });
        expect(wrapper.find('FormDropdown').length).toEqual(7);
      });

      it('should display the docktorGroupURL as readOnly', () => {
        const { wrapper } = setup({ isAdmin: true });
        expect(wrapper.find('[name="docktorGroupURL"]').props().readOnly).toBeFalsy();
      });
    });

    describe('with !isAdmin set to true', () => {
      it('should display the docktorGroupURL as readOnly', () => {
        const { wrapper } = setup({ isAdmin: false });
        expect(wrapper.find('[name="docktorGroupURL"]').props().readOnly).toBeTruthy();
      });
    });

    describe('with no rights management', () => {
      it('should display no <Form.Dropdown />', () => {
        const { wrapper } = setup();
        expect(wrapper.find('FormDropdown').length).toEqual(1);
      });
    });
  });
});

describe('renderTechnologiesField', () => {
  it('renders a drop down with technologies', () => {
    const projectComponent = new ProjectComponent(defaultProps);
    const result = projectComponent.renderMultipleSearchSelectionDropdown(undefined, 'Technologies', ['java', '.NET', 'Pega'], ['java', '.NET', 'Pega'], ['java', '.NET'], 'Java, .NET...', false);
    expect(result).not.toBeNull();
  });
  it('renders div with technologies in readonly mode', () => {
    const projectComponent = new ProjectComponent(defaultProps);
    const result = projectComponent.renderMultipleSearchSelectionDropdown('Technologies', undefined, ['java', '.NET', 'Pega'], ['java', '.NET'], 'Java, .NET...', true);
    expect(result).not.toBeNull();
  });
  it('renders a drop down with empty technologies', () => {
    const projectComponent = new ProjectComponent(defaultProps);
    const result = projectComponent.renderMultipleSearchSelectionDropdown(undefined, 'Technologies', [], ['java', '.NET'], 'Java, .NET...', false);
    expect(result).not.toBeNull();
  });
  it('renders div with empty technologies in readonly mode', () => {
    const projectComponent = new ProjectComponent(defaultProps);
    const result = projectComponent.renderMultipleSearchSelectionDropdown('Technologies', undefined, [], ['java', '.NET'], 'Java, .NET...', true);
    expect(result).not.toBeNull();
  });
  it('renders a drop down with null technologies', () => {
    const projectComponent = new ProjectComponent(defaultProps);
    const result = projectComponent.renderMultipleSearchSelectionDropdown(undefined, 'Technologies', null, ['java', '.NET'], 'Java, .NET...', false);
    expect(result).not.toBeNull();
  });
  it('renders div with null technologies in readonly mode', () => {
    const projectComponent = new ProjectComponent(defaultProps);
    const result = projectComponent.renderMultipleSearchSelectionDropdown('Technologies', undefined, null, ['java', '.NET'], 'Java, .NET...', true);
    expect(result).not.toBeNull();
  });
  it('renders a drop down with undefined technologies', () => {
    const projectComponent = new ProjectComponent(defaultProps);
    const result = projectComponent.renderMultipleSearchSelectionDropdown(undefined, 'Technologies', undefined, ['java', '.NET'], 'Java, .NET...', false);
    expect(result).not.toBeNull();
  });
  it('renders div with undefined technologies in readonly mode', () => {
    const projectComponent = new ProjectComponent(defaultProps);
    const result = projectComponent.renderMultipleSearchSelectionDropdown('Technologies', undefined, undefined, ['java', '.NET'], 'Java, .NET...', true);
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
    const result = projectComponent.renderConsolidationCriteriaField(['Rennes', 'BIOS'], true, null);
    expect(result).not.toBeNull();
  });
  it('renders a drop down with empty consolidation criterias', () => {
    const projectComponent = new ProjectComponent(defaultProps);
    const result = projectComponent.renderConsolidationCriteriaField([], false, null);
    expect(result).not.toBeNull();
  });
  it('renders div with empty consolidation criterias in readonly mode', () => {
    const projectComponent = new ProjectComponent(defaultProps);
    const result = projectComponent.renderConsolidationCriteriaField([], true, null);
    expect(result).not.toBeNull();
  });
  it('renders a drop down with null consolidation criterias', () => {
    const projectComponent = new ProjectComponent(defaultProps);
    const result = projectComponent.renderConsolidationCriteriaField(null, false, null);
    expect(result).not.toBeNull();
  });
  it('renders div with null consolidation criterias in readonly mode', () => {
    const projectComponent = new ProjectComponent(defaultProps);
    const result = projectComponent.renderConsolidationCriteriaField(null, true, null);
    expect(result).not.toBeNull();
  });
  it('renders a drop down with undefined consolidation criterias', () => {
    const projectComponent = new ProjectComponent(defaultProps);
    const result = projectComponent.renderConsolidationCriteriaField(undefined, false, null);
    expect(result).not.toBeNull();
  });
  it('renders div with undefined consolidation criterias in readonly mode', () => {
    const projectComponent = new ProjectComponent(defaultProps);
    const result = projectComponent.renderConsolidationCriteriaField(undefined, true, null);
    expect(result).not.toBeNull();
  });
});
