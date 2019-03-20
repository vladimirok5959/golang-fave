/// <reference types="Cypress" />

context('Install MySQL, create first user and login', () => {
  it('should do redirect to cp panel', () => {
    cy.resetCMS();
    cy.request({
      url: cy.getBaseUrl() + '/',
      followRedirect: false
    }).then((response) => {
      expect(response.status).to.eq(302);
      expect(response.redirectedToUrl).to.eq(cy.getBaseUrl() + '/cp/');
    });
    cy.visitCMS('/cp/');
    cy.url().should('eq', cy.getBaseUrl() + '/cp/');
  });

  it('should configure mysql config', () => {
    cy.actionStart();
    cy.get('.form-signin input[type=text]').should('have.length', 4);
    cy.get('.form-signin input[type=password]').should('have.length', 1);
    cy.get('.form-signin button').should('have.length', 1);
    cy.get('.form-signin input[name=name]').type('fave');
    cy.get('.form-signin input[name=user]').type('root');
    cy.get('.form-signin input[name=password]').type('root');
    cy.get('.form-signin button').click();
    cy.actionWait();
  });

  it('should create first user', () => {
    cy.actionStart();
    cy.get('.form-signin input[type=text]').should('have.length', 2);
    cy.get('.form-signin input[type=email]').should('have.length', 1);
    cy.get('.form-signin input[type=password]').should('have.length', 1);
    cy.get('.form-signin button').should('have.length', 1);
    cy.get('.form-signin input[name=first_name]').type('First');
    cy.get('.form-signin input[name=last_name]').type('Last');
    cy.get('.form-signin input[name=email]').type('example@example.com');
    cy.get('.form-signin input[name=password]').type('example@example.com');
    cy.get('.form-signin button').click();
    cy.actionWait();
  });

  it('should login to control panel', () => {
    cy.actionStart();
    cy.get('.form-signin input[type=email]').should('have.length', 1);
    cy.get('.form-signin input[type=password]').should('have.length', 1);
    cy.get('.form-signin button').should('have.length', 1);
    cy.get('.form-signin input[name=email]').type('example@example.com');
    cy.get('.form-signin input[name=password]').type('example@example.com');
    cy.get('.form-signin button').click();
    cy.actionWait();
    cy.logoutCMS();
  });
});
