// import constants
import ModalConstants from './modal.constants';

// Close Modal
const closeModal = () => {
  return {
    type: ModalConstants.CLOSE_MODAL
  };
};

// New Url Modal
const openNewUrlModal = (callback) => {
  let form = { lines: [], hidden: [] };
  let line = { class: 'two', fields: [] };
  line.fields.push({ label: 'Name', name: 'name', placeholder: 'URL Name', type: 'text', required: true });
  line.fields.push({ label: 'Link', name: 'link', placeholder: 'URL Link', type: 'url', required: true });
  form.lines.push(line);
  return {
    type: ModalConstants.OPEN_MODAL,
    title: 'New Url',
    form,
    callback
  };
};

// Edit Url Modal
const openEditUrlModal = (url, callback) => {
  let form = { lines: [], hidden: [] };
  let line = { class: 'two', fields: [] };
  line.fields.push({ label: 'Name', name: 'name', placeholder: 'URL Name', type: 'text', required: true, value: url.name });
  line.fields.push({ label: 'Link', name: 'link', placeholder: 'URL Link', type: 'url', required: true, value: url.link });
  form.lines.push(line);
  return {
    type: ModalConstants.OPEN_MODAL,
    title: 'Edit Url',
    form,
    callback
  };
};

const openRemoveProjectModal = (project, callback) => {
  let form = { lines: [], hidden: [] };
  return {
    type: ModalConstants.OPEN_MODAL,
    title: `${project.name} - Removing project`,
    basic: true,
    message: 'This project will be removed for all users. Are you sure to delete it?',
    form,
    callback
  };
};

export default {
  closeModal,
  openNewUrlModal,
  openEditUrlModal,
  openRemoveProjectModal
};
