/// <reference types="Cypress" />

context('Module users', () => {
  it('should reset', () => {
    cy.installCMS();
  });

  it('should render data table', () => {
    cy.loginCMS();
    cy.visitCMS('/cp/users/');
    cy.get('table.data-table thead tr').should('have.length', 1);
    cy.get('table.data-table thead tr th').should('have.length', 4);
    cy.get('table.data-table tbody tr').should('have.length', 1);
    cy.get('table.data-table tbody tr:nth-child(1) td').should('have.length', 4);
    cy.logoutCMS();
  });

  it('should render data form', () => {
    cy.loginCMS();
    cy.visitCMS('/cp/users/add/');
    cy.get('.data-form.users-add input[type=text]').should('have.length', 2);
    cy.get('.data-form.users-add input[type=email]').should('have.length', 1);
    cy.get('.data-form.users-add input[type=password]').should('have.length', 1);
    cy.get('.data-form.users-add input[type=checkbox]').should('have.length', 2);
    cy.logoutCMS();
  });

  it('should not add new user', () => {
    cy.loginCMS();
    cy.visitCMS('/cp/users/add/');
    cy.get('.data-form.users-add input[name=email]').clear().type('some@text');
    cy.get('.data-form.users-add input[name=password]').clear().type('some@text');
    cy.get('#add-edit-button').click();
    cy.actionWait();
    cy.get('#sys-modal-system-message').should('exist');
    cy.get('#sys-modal-system-message .modal-body').contains('Please specify correct user email');
    cy.get('#sys-modal-system-message .modal-footer').find('button').click();
    cy.logoutCMS();
  });

  it('should add new user', () => {
    cy.loginCMS();
    cy.visitCMS('/cp/users/add/');
    cy.get('.data-form.users-add input[name=first_name]').clear().type('Some user first name');
    cy.get('.data-form.users-add input[name=last_name]').clear().type('Some user last name');
    cy.get('.data-form.users-add input[name=email]').clear().type('some@user.com');
    cy.get('.data-form.users-add input[name=password]').clear().type('some@text');
    cy.get('.data-form.users-add label[for=lbl_active]').click();
    cy.get('.data-form.users-add label[for=lbl_admin]').click();
    cy.get('#add-edit-button').click();
    cy.actionWait();
    cy.logoutCMS();
  });

  it('should render added user in list', () => {
    cy.loginCMS();
    cy.visitCMS('/cp/users/');
    cy.get('table.data-table tbody tr').should('have.length', 2);
    cy.get('table.data-table tbody tr td').should('contain', 'some@user.com');
    cy.contains('table.data-table tbody tr td a', 'some@user.com').parentsUntil('tr').parent().find('.svg-green').should('exist');
    cy.logoutCMS();
  });

  it('should render added user in edit form', () => {
    cy.loginCMS();
    cy.visitCMS('/cp/users/');
    cy.contains('table.data-table tbody tr td a', 'some@user.com').click();
    cy.get('.data-form.users-modify input[name=first_name]').should('have.value', 'Some user first name');
    cy.get('.data-form.users-modify input[name=last_name]').should('have.value', 'Some user last name');
    cy.get('.data-form.users-modify input[name=email]').should('have.value', 'some@user.com');
    cy.get('.data-form.users-modify input[name=password]').should('have.value', '');
    cy.get('.data-form.users-modify input[name=active]').should('be.checked');
    cy.get('.data-form.users-modify input[name=admin]').should('be.checked');
    cy.logoutCMS();
  });

  it('should delete added user', () => {
    cy.loginCMS();
    cy.visitCMS('/cp/users/');
    cy.contains('table.data-table tbody tr td a', 'some@user.com').parentsUntil('tr').parent().find('td a.ico.delete').click();
    cy.actionWait();
    cy.get('table.data-table tbody tr').should('have.length', 1);
    cy.logoutCMS();
  });
});
