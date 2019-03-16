/// <reference types="Cypress" />

context('Module pages', () => {
  it('should reset', () => {
    cy.installCMS();
  });

  it('should render data table', () => {
    cy.loginCMS();
    cy.visit('http://localhost:8080/cp/');
    cy.get('table#cp-table-pages thead tr').should('have.length', 1);
    cy.get('table#cp-table-pages thead tr th').should('have.length', 4);
    cy.get('table#cp-table-pages tbody tr').should('have.length', 3);
    cy.get('table#cp-table-pages tbody tr:nth-child(1) td').should('have.length', 4);
  });

  it('should render data form', () => {
    cy.loginCMS();
    cy.visit('http://localhost:8080/cp/index/add/');
    
    cy.get('.data-form.index-add input[type=text]').should('have.length', 4);
    cy.get('.data-form.index-add textarea').should('have.length', 2);
    cy.get('.data-form.index-add input[type=checkbox]').should('have.length', 1);
  });

  it('should not add new page', () => {
    cy.loginCMS();
    cy.visit('http://localhost:8080/cp/index/add/');
    cy.get('#add-edit-button').click();
    cy.actionWait();

    cy.get('.data-form.index-add div.sys-messages').should('exist');
  });

  it('should add new page', () => {
    cy.loginCMS();
    cy.visit('http://localhost:8080/cp/index/add/');
    
    cy.get('.data-form.index-add input[name=name]').clear().type('Some test page');
    cy.get('.data-form.index-add textarea[name=content]').clear().type('Some test content');
    cy.get('.data-form.index-add input[name=meta_title]').clear().type('Page meta title');
    cy.get('.data-form.index-add input[name=meta_keywords]').clear().type('Page meta keywords');
    cy.get('.data-form.index-add textarea[name=meta_description]').clear().type('Page meta description');
    cy.get('.data-form.index-add label[for=lbl_active]').click();
    cy.get('#add-edit-button').click();
    cy.actionWait();
  });

  it('should render added page', () => {
    cy.loginCMS();
    cy.visit('http://localhost:8080/cp/');
    cy.get('table#cp-table-pages tbody tr').should('have.length', 4);
    cy.get('table#cp-table-pages tbody tr:nth-child(1) td:nth-child(1)').should('contain', 'Some test page');
    cy.get('table#cp-table-pages tbody tr:nth-child(1) td:nth-child(3) .svg-green').should('exist');
  });

  it('should delete added page', () => {
    cy.loginCMS();
    cy.visit('http://localhost:8080/cp/');
    cy.get('table#cp-table-pages tbody tr:nth-child(1) td:nth-child(4) a.ico.delete').click();
    cy.actionWait();
    cy.get('table#cp-table-pages tbody tr').should('have.length', 3);
  });
});
