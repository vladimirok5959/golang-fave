/// <reference types="Cypress" />

context('Module blog categories', () => {
  it('should reset', () => {
    cy.installCMS();
  });

  it('should render data table', () => {
    cy.loginCMS();
    cy.visit('http://localhost:8080/cp/blog/categories/');
    cy.get('table.data-table thead tr').should('have.length', 1);
    cy.get('table.data-table thead tr th').should('have.length', 2);
    cy.get('table.data-table tbody tr').should('have.length', 11);
    cy.get('table.data-table tbody tr:nth-child(1) td').should('have.length', 2);
    cy.logoutCMS();
  });

  it('should render data form', () => {
    cy.loginCMS();
    cy.visit('http://localhost:8080/cp/blog/categories-add/');
    cy.get('.data-form.blog-categories-add select').should('have.length', 1);
    cy.get('.data-form.blog-categories-add input[type=text]').should('have.length', 2);
    cy.logoutCMS();
  });

  it('should not add new category', () => {
    cy.loginCMS();
    cy.visit('http://localhost:8080/cp/blog/categories-add/');
    cy.get('#add-edit-button').click();
    cy.actionWait();
    cy.get('.data-form.blog-categories-add div.sys-messages').should('exist');
    cy.logoutCMS();
  });

  it('should add new category', () => {
    cy.loginCMS();
    cy.visit('http://localhost:8080/cp/blog/categories-add/');
    cy.get('.data-form.blog-categories-add input[name=name]').clear().type('Some test category');
    cy.get('#add-edit-button').click();
    cy.actionWait();
    cy.logoutCMS();
  });

  it('should render added category in list', () => {
    cy.loginCMS();
    cy.visit('http://localhost:8080/cp/blog/categories/');
    cy.get('table.data-table tbody tr').should('have.length', 12);
    cy.get('table.data-table tbody tr td').should('contain', 'Some test category');
    cy.logoutCMS();
  });

  it('should render added category in edit form', () => {
    cy.loginCMS();
    cy.visit('http://localhost:8080/cp/blog/categories/');
    cy.contains('table.data-table tbody tr td a', 'Some test category').click();
    cy.get('.data-form.blog-categories-modify select[name=parent]').should('have.value', '0');
    cy.get('.data-form.blog-categories-modify input[name=name]').should('have.value', 'Some test category');
    cy.get('.data-form.blog-categories-modify input[name=alias]').should('have.value', 'some-test-category');
    cy.logoutCMS();
  });

  it('should add new child category', () => {
    cy.loginCMS();
    cy.visit('http://localhost:8080/cp/blog/categories-add/');
    cy.get('.data-form.blog-categories-add select[name=parent]').select('Some test category');
    cy.get('.data-form.blog-categories-add input[name=name]').clear().type('Some test child category');
    cy.get('#add-edit-button').click();
    cy.actionWait();
    cy.logoutCMS();
  });

  it('should render added child category in list', () => {
    cy.loginCMS();
    cy.visit('http://localhost:8080/cp/blog/categories/');
    cy.get('table.data-table tbody tr').should('have.length', 13);
    cy.get('table.data-table tbody tr td').should('contain', '— Some test child category');
    cy.logoutCMS();
  });

  it('should render added child category in edit form', () => {
    cy.loginCMS();
    cy.visit('http://localhost:8080/cp/blog/categories/');
    cy.contains('table.data-table tbody tr td a', '— Some test child category').click();
    cy.get('.data-form.blog-categories-modify select[name=parent]').find(':selected').contains('Some test category')
    cy.get('.data-form.blog-categories-modify input[name=name]').should('have.value', 'Some test child category');
    cy.get('.data-form.blog-categories-modify input[name=alias]').should('have.value', 'some-test-child-category');
    cy.logoutCMS();
  });

  it('should delete added child category', () => {
    cy.loginCMS();
    cy.visit('http://localhost:8080/cp/blog/categories/');
    cy.contains('table.data-table tbody tr td a', '— Some test child category').parentsUntil('tr').parent().find('td a.ico.delete').click();
    cy.actionWait();
    cy.get('table.data-table tbody tr').should('have.length', 12);
    cy.logoutCMS();
  });

  it('should delete added category', () => {
    cy.loginCMS();
    cy.visit('http://localhost:8080/cp/blog/categories/');
    cy.contains('table.data-table tbody tr td a', 'Some test category').parentsUntil('tr').parent().find('td a.ico.delete').click();
    cy.actionWait();
    cy.get('table.data-table tbody tr').should('have.length', 11);
    cy.logoutCMS();
  });
});
