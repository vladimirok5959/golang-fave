/// <reference types="Cypress" />

context('Module pages', () => {
  it('should reset', () => {
    cy.installCMS();
  });

  it('should render data table', () => {
    cy.loginCMS();
    cy.visitCMS('/cp/');
    cy.get('table.data-table thead tr').should('have.length', 1);
    cy.get('table.data-table thead tr th').should('have.length', 4);
    cy.get('table.data-table tbody tr').should('have.length', 3);
    cy.get('table.data-table tbody tr:nth-child(1) td').should('have.length', 4);
    cy.logoutCMS();
  });

  it('should render data form', () => {
    cy.loginCMS();
    cy.visitCMS('/cp/index/add/');
    cy.get('.data-form.index-add input[type=text]').should('have.length', 4);
    cy.get('.data-form.index-add textarea').should('have.length', 2);
    cy.get('.data-form.index-add input[type=checkbox]').should('have.length', 1);
    cy.logoutCMS();
  });

  it('should not add new page', () => {
    cy.loginCMS();
    cy.visitCMS('/cp/index/add/');
    cy.get('#add-edit-button').click();
    cy.actionWait();
    cy.get('.data-form.index-add div.sys-messages').should('exist');
    cy.logoutCMS();
  });

  it('should add new page', () => {
    cy.loginCMS();
    cy.visitCMS('/cp/index/add/');
    cy.get('.data-form.index-add input[name=name]').clear().type('Some test page');
    cy.get('.data-form.index-add textarea[name=content]').parent().find('.pell-content').clear().type('Some test content');
    cy.get('.data-form.index-add input[name=meta_title]').clear().type('Page meta title');
    cy.get('.data-form.index-add input[name=meta_keywords]').clear().type('Page meta keywords');
    cy.get('.data-form.index-add textarea[name=meta_description]').clear().type('Page meta description');
    cy.get('.data-form.index-add label[for=lbl_active]').click();
    cy.get('#add-edit-button').click();
    cy.actionWait();
    cy.logoutCMS();
  });

  it('should render added page in list', () => {
    cy.loginCMS();
    cy.visitCMS('/cp/');
    cy.get('table.data-table tbody tr').should('have.length', 4);
    cy.get('table.data-table tbody tr td').should('contain', 'Some test page');
    cy.contains('table#cp-table-pages tbody tr td a', 'Some test page').parentsUntil('tr').parent().find('.svg-green').should('exist');
    cy.logoutCMS();
  });

  it('should render added page in edit form', () => {
    cy.loginCMS();
    cy.visitCMS('/cp/');
    cy.contains('table.data-table tbody tr td a', 'Some test page').click();
    cy.get('.data-form.index-modify input[name=name]').should('have.value', 'Some test page');
    cy.get('.data-form.index-modify input[name=alias]').should('have.value', '/some-test-page/');
    cy.get('.data-form.index-modify textarea[name=content]').parent().find('.pell-content').should(($editor) => {
      expect($editor).to.have.text('Some test content');
    });
    cy.get('.data-form.index-modify input[name=meta_title]').should('have.value', 'Page meta title');
    cy.get('.data-form.index-modify input[name=meta_keywords]').should('have.value', 'Page meta keywords');
    cy.get('.data-form.index-modify textarea[name=meta_description]').should('have.value', 'Page meta description');
    cy.get('.data-form.index-modify input[name=active]').should('be.checked');
    cy.logoutCMS();
  });

  it('should delete added page', () => {
    cy.loginCMS();
    cy.visitCMS('/cp/');
    cy.contains('table.data-table tbody tr td a', 'Some test page').parentsUntil('tr').parent().find('td a.ico.delete').click();
    cy.actionWait();
    cy.get('table.data-table tbody tr').should('have.length', 3);
    cy.logoutCMS();
  });
});
