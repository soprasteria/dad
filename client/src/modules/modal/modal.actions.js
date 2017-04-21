// import constants
import ModalConstants from './modal.constants';

// Close Modal
const closeModal = () => {
  return {
    type: ModalConstants.CLOSE_MODAL
  };
};

const openRemoveProjectModal = (project, callback) => {
  let form = { lines: [], hidden: [] };
  return {
    type: ModalConstants.OPEN_MODAL,
    icon: 'trash',
    title: `${project.name} - Removing project`,
    message: 'This project will be removed for all users. Are you sure to delete it?',
    validateText: 'Remove',
    validateColor: 'red',
    validateIcon: 'trash',
    cancelText:'No',
    cancelColor: 'black',
    form,
    callback
  };
};

export default {
  closeModal,
  openRemoveProjectModal
};
