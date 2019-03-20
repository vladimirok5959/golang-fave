/// <reference types="Cypress" />

context('Module blog posts', () => {
  it('should reset', () => {
    cy.installCMS();
  });

  it('should render data table', () => {
    cy.loginCMS();
    cy.visitCMS('/cp/blog/');
    cy.get('table.data-table thead tr').should('have.length', 1);
    cy.get('table.data-table thead tr th').should('have.length', 4);
    cy.get('table.data-table tbody tr').should('have.length', 3);
    cy.get('table.data-table tbody tr:nth-child(1) td').should('have.length', 4);
    cy.logoutCMS();
  });

  it('should render data form', () => {
    cy.loginCMS();
    cy.visitCMS('/cp/blog/add/');
    cy.get('.data-form.blog-add input[type=text]').should('have.length', 2);
    cy.get('.data-form.blog-add select').should('have.length', 1);
    cy.get('.data-form.blog-add textarea').should('have.length', 1);
    cy.get('.data-form.blog-add input[type=checkbox]').should('have.length', 1);
    cy.logoutCMS();
  });

  it('should not add new post', () => {
    cy.loginCMS();
    cy.visitCMS('/cp/blog/add/');
    cy.get('#add-edit-button').click();
    cy.actionWait();
    cy.get('.data-form.blog-add div.sys-messages').should('exist');
    cy.logoutCMS();
  });

  it('should add new post', () => {
    cy.loginCMS();
    cy.visitCMS('/cp/blog/add/');
    cy.get('.data-form.blog-add input[name=name]').clear().type('Some test post');
    cy.get('.data-form.blog-add select#lbl_cats').select(['Health and food', '— — Natural']).invoke('val').should('deep.equal', ['2', '7']);
    cy.get('.data-form.blog-add textarea[name=content]').clear().type('Some test content');
    cy.get('.data-form.blog-add label[for=lbl_active]').click();
    cy.get('#add-edit-button').click();
    cy.actionWait();
    cy.logoutCMS();
  });

  it('should render added post in list', () => {
    cy.loginCMS();
    cy.visitCMS('/cp/blog/');
    cy.get('table.data-table tbody tr').should('have.length', 4);
    cy.get('table.data-table tbody tr td').should('contain', 'Some test post');
    cy.contains('table.data-table tbody tr td a', 'Some test post').parentsUntil('tr').parent().find('.svg-green').should('exist');
    cy.logoutCMS();
  });

  it('should render added post in edit form', () => {
    cy.loginCMS();
    cy.visitCMS('/cp/blog/');
    cy.contains('table.data-table tbody tr td a', 'Some test post').click();
    cy.get('.data-form.blog-modify input[name=name]').should('have.value', 'Some test post');
    cy.get('.data-form.blog-modify input[name=alias]').should('have.value', 'some-test-post');
    cy.get('.data-form.blog-modify select#lbl_cats').invoke('val').should('deep.equal', ['2', '7']);
    cy.get('.data-form.blog-modify textarea[name=content]').should('have.value', 'Some test content');
    cy.get('.data-form.blog-modify input[name=active]').should('be.checked');
    cy.logoutCMS();
  });

  it('should delete added post', () => {
    cy.loginCMS();
    cy.visitCMS('/cp/blog/');
    cy.contains('table.data-table tbody tr td a', 'Some test post').parentsUntil('tr').parent().find('td a.ico.delete').click();
    cy.actionWait();
    cy.get('table.data-table tbody tr').should('have.length', 3);
    cy.logoutCMS();
  });
});
