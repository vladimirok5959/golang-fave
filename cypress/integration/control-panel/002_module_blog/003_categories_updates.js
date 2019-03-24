/// <reference types="Cypress" />

context('Module blog categories updates', () => {
  it('should reset', () => {
    cy.installCMS();
  });

  it('should render correct data in table', () => {
    cy.loginCMS();
    cy.visitCMS('/cp/blog/categories/');
    cy.get('table.data-table thead tr').should('have.length', 1);
    cy.get('table.data-table thead tr th').should('have.length', 2);
    cy.get('table.data-table tbody tr').should('have.length', 11);
    cy.get('table.data-table tbody tr:nth-child(1) td').should('have.length', 2);

    cy.get('table.data-table tbody tr:nth-child(1) td:nth-child(1)').should('contain', 'Health and food');
    cy.get('table.data-table tbody tr:nth-child(2) td:nth-child(1)').should('contain', '— Juices');
    cy.get('table.data-table tbody tr:nth-child(3) td:nth-child(1)').should('contain', '— — Natural');
    cy.get('table.data-table tbody tr:nth-child(4) td:nth-child(1)').should('contain', '— — For kids');
    cy.get('table.data-table tbody tr:nth-child(5) td:nth-child(1)').should('contain', '— Nutrition');
    cy.get('table.data-table tbody tr:nth-child(6) td:nth-child(1)').should('contain', '— — For all');
    cy.get('table.data-table tbody tr:nth-child(7) td:nth-child(1)').should('contain', '— — For athletes');
    cy.get('table.data-table tbody tr:nth-child(8) td:nth-child(1)').should('contain', 'News');
    cy.get('table.data-table tbody tr:nth-child(9) td:nth-child(1)').should('contain', '— Computers and technology');
    cy.get('table.data-table tbody tr:nth-child(10) td:nth-child(1)').should('contain', '— Film industry');
    cy.get('table.data-table tbody tr:nth-child(11) td:nth-child(1)').should('contain', 'Hobby');

    cy.logoutCMS();
  });

  it('should change category parent (from left to right)', () => {
    cy.loginCMS();
    cy.visitCMS('/cp/blog/categories/');
    
    cy.contains('table.data-table tbody tr td a', '— Juices').click();
    cy.get('.data-form.blog-categories-modify select[name=parent]').select('News');
    cy.get('#add-edit-button').click();
    cy.actionWait();

    cy.visitCMS('/cp/blog/categories/');

    cy.get('table.data-table tbody tr:nth-child(1) td:nth-child(1)').should('contain', 'Health and food');
    cy.get('table.data-table tbody tr:nth-child(2) td:nth-child(1)').should('contain', '— Nutrition');
    cy.get('table.data-table tbody tr:nth-child(3) td:nth-child(1)').should('contain', '— — For all');
    cy.get('table.data-table tbody tr:nth-child(4) td:nth-child(1)').should('contain', '— — For athletes');
    cy.get('table.data-table tbody tr:nth-child(5) td:nth-child(1)').should('contain', 'News');
    cy.get('table.data-table tbody tr:nth-child(6) td:nth-child(1)').should('contain', '— Computers and technology');
    cy.get('table.data-table tbody tr:nth-child(7) td:nth-child(1)').should('contain', '— Film industry');
    cy.get('table.data-table tbody tr:nth-child(8) td:nth-child(1)').should('contain', '— Juices');
    cy.get('table.data-table tbody tr:nth-child(9) td:nth-child(1)').should('contain', '— — Natural');
    cy.get('table.data-table tbody tr:nth-child(10) td:nth-child(1)').should('contain', '— — For kids');
    cy.get('table.data-table tbody tr:nth-child(11) td:nth-child(1)').should('contain', 'Hobby');

    cy.logoutCMS();
  });

  it('should change category parent (from right to left)', () => {
    cy.loginCMS();
    cy.visitCMS('/cp/blog/categories/');
    
    cy.contains('table.data-table tbody tr td a', '— Juices').click();
    cy.get('.data-form.blog-categories-modify select[name=parent]').select('— Nutrition');
    cy.get('#add-edit-button').click();
    cy.actionWait();

    cy.visitCMS('/cp/blog/categories/');

    cy.get('table.data-table tbody tr:nth-child(1) td:nth-child(1)').should('contain', 'Health and food');
    cy.get('table.data-table tbody tr:nth-child(2) td:nth-child(1)').should('contain', '— Nutrition');
    cy.get('table.data-table tbody tr:nth-child(3) td:nth-child(1)').should('contain', '— — For all');
    cy.get('table.data-table tbody tr:nth-child(4) td:nth-child(1)').should('contain', '— — For athletes');
    cy.get('table.data-table tbody tr:nth-child(5) td:nth-child(1)').should('contain', '— — Juices');
    cy.get('table.data-table tbody tr:nth-child(6) td:nth-child(1)').should('contain', '— — — Natural');
    cy.get('table.data-table tbody tr:nth-child(7) td:nth-child(1)').should('contain', '— — — For kids');
    cy.get('table.data-table tbody tr:nth-child(8) td:nth-child(1)').should('contain', 'News');
    cy.get('table.data-table tbody tr:nth-child(9) td:nth-child(1)').should('contain', '— Computers and technology');
    cy.get('table.data-table tbody tr:nth-child(10) td:nth-child(1)').should('contain', '— Film industry');
    cy.get('table.data-table tbody tr:nth-child(11) td:nth-child(1)').should('contain', 'Hobby');

    cy.logoutCMS();
  });

  it('should do not allow to change category parent to they child as parent', () => {
    cy.loginCMS();
    cy.visitCMS('/cp/blog/categories/');

    cy.contains('table.data-table tbody tr td a', '— Juices').click();
    cy.get('.data-form.blog-categories-modify select[name=parent]').select('— — — Natural');
    cy.get('#add-edit-button').click();
    cy.actionWait();

    cy.get('.data-form.blog-categories-modify div.sys-messages').should('exist');

    cy.logoutCMS();
  });
});
