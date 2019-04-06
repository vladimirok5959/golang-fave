/// <reference types="Cypress" />

context('Module Pagination', () => {
  it('should reset', () => {
    cy.installCMS();
  });

  it('should render inputs for blog pagination', () => {
    cy.loginCMS();
    cy.visitCMS('/cp/settings/pagination/');
    cy.get('.data-form.settings-pagination input[name=blog-index]').should('exist');
    cy.get('.data-form.settings-pagination input[name=blog-category]').should('exist');
    cy.get('.data-form.settings-pagination input[name=blog-index]').should('have.value', '5');
    cy.get('.data-form.settings-pagination input[name=blog-category]').should('have.value', '5');
    cy.logoutCMS();
  });

  it('should change inputs value for blog pagination', () => {
    cy.loginCMS();
    cy.visitCMS('/cp/settings/pagination/');
    cy.get('.data-form.settings-pagination input[name=blog-index]').clear().type('2');
    cy.get('.data-form.settings-pagination input[name=blog-category]').clear().type('3');
    cy.get('#add-edit-button').click();
    cy.actionWait();
    cy.visitCMS('/cp/settings/pagination/');
    cy.get('.data-form.settings-pagination input[name=blog-index]').should('have.value', '2');
    cy.get('.data-form.settings-pagination input[name=blog-category]').should('have.value', '3');
    cy.get('.data-form.settings-pagination input[name=blog-index]').clear().type('5');
    cy.get('.data-form.settings-pagination input[name=blog-category]').clear().type('5');
    cy.get('#add-edit-button').click();
    cy.actionWait();
    cy.logoutCMS();
  });
});
