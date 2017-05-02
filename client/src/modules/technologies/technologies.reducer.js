import { generateEntitiesReducer } from '../utils/entities';

const technologiesReducer = (state, action) => {
  const entitiesState = generateEntitiesReducer(state, action, 'technologies');
  switch (action.type) {
  default:
    return entitiesState;
  }
};

export default technologiesReducer;
