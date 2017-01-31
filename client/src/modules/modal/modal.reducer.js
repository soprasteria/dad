// import constants
import ModalConstants from './modal.constants';

const initialState = {
  isVisible: false,
  form: { lines: [], hidden: [] }
};

const modalReducer = (state = initialState, action) => {
  switch (action.type) {
  case ModalConstants.CLOSE_MODAL:
    return { ...initialState };

  case ModalConstants.OPEN_MODAL:
    return { ...initModal(action) };

  default:
    return state;
  }
};

const initModal = (action) => {
  let res = {};
  res.isVisible = true;
  res.title = action.title;
  res.form = action.form;
  res.callback = action.callback;
  res.basic = action.basic || false;
  res.message = action.message || '';
  res.icon = action.icon;
  res.validateText = action.validateText || 'Validate';
  res.validateColor = action.validateColor || 'green';
  res.validateIcon = action.validateIcon || 'checkmark';
  res.cancelText = action.cancelText || 'Cancel';
  res.cancelColor = action.cancelColor || 'black';
  return res;
};

export default modalReducer;
