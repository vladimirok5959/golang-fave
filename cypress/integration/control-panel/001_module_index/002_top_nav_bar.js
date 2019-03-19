/// <reference types="Cypress" />

context('Top navigation bar', () => {
  it('should reset', () => {
    cy.installCMS();
  });

  it('should render top nav bar', () => {
    cy.loginCMS();
    cy.get('#navbarCollapse ul.navbar-nav').should('have.length', 2);
    cy.logoutCMS();
  });

  it('should render modules menu', () => {
    cy.loginCMS();
    cy.get('#navbarCollapse ul.navbar-nav:nth-child(1) li.nav-item').should('have.length', 2);
    cy.get('#navbarCollapse ul.navbar-nav:nth-child(1) li.nav-item:nth-child(1) a.nav-link').should('exist');
    cy.get('#navbarCollapse ul.navbar-nav:nth-child(1) li.nav-item:nth-child(2) a.nav-link').should('exist');
    cy.logoutCMS();
  });

  it('should render user menu', () => {
    cy.loginCMS();
    cy.get('#navbarCollapse ul.navbar-nav:nth-child(2) li.nav-item').should('have.length', 1);
    cy.get('#navbarCollapse ul.navbar-nav:nth-child(2) li.nav-item:nth-child(1) a.nav-link').should('contain', 'example@example.com');
    cy.logoutCMS();
  });

  it('should render user profile modal dialog', () => {
    cy.loginCMS();
    cy.get('#navbarCollapse ul.navbar-nav:nth-child(2) li.nav-item:nth-child(1) a.nav-link:first-child').click();
    cy.get('#navbarCollapse ul.navbar-nav:nth-child(2) li.nav-item:nth-child(1) div.dropdown-menu a.dropdown-item:nth-child(1)').click();
    cy.get('#sys-modal-user-settings').should('exist');
    cy.get('#sys-modal-user-settings form input[type=text]').should('have.length', 2);
    cy.get('#sys-modal-user-settings form input[type=email]').should('have.length', 1);
    cy.get('#sys-modal-user-settings form input[type=password]').should('have.length', 1);
    cy.logoutCMS();
  });

  it('should change user profile data in modal dialog', () => {
    cy.loginCMS();
    cy.get('#navbarCollapse ul.navbar-nav:nth-child(2) li.nav-item:nth-child(1) a.nav-link:first-child').click();
    cy.get('#navbarCollapse ul.navbar-nav:nth-child(2) li.nav-item:nth-child(1) div.dropdown-menu a.dropdown-item:nth-child(1)').click();
    cy.get('#sys-modal-user-settings').should('exist');
    cy.get('#sys-modal-user-settings form input[name=first_name]').clear().type('FirstNew');
    cy.get('#sys-modal-user-settings form input[name=last_name]').clear().type('LastNew');
    cy.get('#sys-modal-user-settings form button[type=submit].btn').click();
    cy.actionWait();
    cy.logoutCMS();
  });

  it('should render saved user profile data in modal dialog', () => {
    cy.loginCMS();
    cy.get('#navbarCollapse ul.navbar-nav:nth-child(2) li.nav-item:nth-child(1) a.nav-link:first-child').click();
    cy.get('#navbarCollapse ul.navbar-nav:nth-child(2) li.nav-item:nth-child(1) div.dropdown-menu a.dropdown-item:nth-child(1)').click();
    cy.get('#sys-modal-user-settings').should('exist');
    cy.get('#sys-modal-user-settings form input[name=first_name]').should('have.value', 'FirstNew');
    cy.get('#sys-modal-user-settings form input[name=last_name]').should('have.value', 'LastNew');
    cy.get('#sys-modal-user-settings form button[type=button].btn').click();
    cy.logoutCMS();
  });
});
