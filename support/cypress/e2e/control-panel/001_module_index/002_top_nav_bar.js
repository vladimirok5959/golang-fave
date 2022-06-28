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
    cy.get('#navbarCollapse ul.navbar-nav:nth-child(1) li.nav-item').should('have.length', 3);
    cy.get('#navbarCollapse ul.navbar-nav li.nav-item a.nav-link').should('contain', 'Modules');
    cy.contains('#navbarCollapse ul.navbar-nav li.nav-item a.nav-link', 'Modules').parent().find('.dropdown-menu .dropdown-item').should('contain', 'Pages');
    cy.contains('#navbarCollapse ul.navbar-nav li.nav-item a.nav-link', 'Modules').parent().find('.dropdown-menu .dropdown-item').should('contain', 'Blog');
    cy.get('#navbarCollapse ul.navbar-nav li.nav-item a.nav-link').should('contain', 'System');
    cy.contains('#navbarCollapse ul.navbar-nav li.nav-item a.nav-link', 'System').parent().find('.dropdown-menu .dropdown-item').should('contain', 'Users');
    cy.contains('#navbarCollapse ul.navbar-nav li.nav-item a.nav-link', 'System').parent().find('.dropdown-menu .dropdown-item').should('contain', 'Settings');
    cy.logoutCMS();
  });

  it('should render user menu', () => {
    cy.loginCMS();
    cy.get('#navbarCollapse ul.navbar-nav:nth-child(2) li.nav-item').should('have.length', 1);
    cy.get('#navbarCollapse ul.navbar-nav li.nav-item a.nav-link').should('contain', 'example@example.com');
    cy.contains('#navbarCollapse ul.navbar-nav li.nav-item a.nav-link', 'example@example.com').parent().find('.dropdown-menu .dropdown-item').should('contain', 'My profile');
    cy.contains('#navbarCollapse ul.navbar-nav li.nav-item a.nav-link', 'example@example.com').parent().find('.dropdown-menu .dropdown-item').should('contain', 'Logout');
    cy.logoutCMS();
  });

  it('should render user profile modal dialog', () => {
    cy.loginCMS();
    cy.contains('#navbarCollapse ul.navbar-nav li.nav-item a.nav-link', 'example@example.com').click();
    cy.contains('#navbarCollapse ul.navbar-nav li.nav-item a.nav-link', 'example@example.com').parent().find('.dropdown-menu').contains('a.dropdown-item', 'My profile').click();
    cy.get('#sys-modal-user-settings').should('exist');
    cy.get('#sys-modal-user-settings form input[type=text]').should('have.length', 2);
    cy.get('#sys-modal-user-settings form input[type=email]').should('have.length', 1);
    cy.get('#sys-modal-user-settings form input[type=password]').should('have.length', 1);
    cy.logoutCMS();
  });

  it('should change user profile data in modal dialog', () => {
    cy.loginCMS();
    cy.contains('#navbarCollapse ul.navbar-nav li.nav-item a.nav-link', 'example@example.com').click();
    cy.contains('#navbarCollapse ul.navbar-nav li.nav-item a.nav-link', 'example@example.com').parent().find('.dropdown-menu').contains('a.dropdown-item', 'My profile').click();
    cy.get('#sys-modal-user-settings').should('exist');
    cy.get('#sys-modal-user-settings form input[name=first_name]').clear();
    cy.wait(500);
    cy.get('#sys-modal-user-settings form input[name=last_name]').clear();
    cy.wait(500);
    cy.get('#sys-modal-user-settings form input[name=first_name]').clear().type('FirstNew');
    cy.get('#sys-modal-user-settings form input[name=last_name]').clear().type('LastNew');
    cy.get('#sys-modal-user-settings form button[type=submit].btn').click();
    cy.actionWait();
    cy.logoutCMS();
  });

  it('should render saved user profile data in modal dialog', () => {
    cy.loginCMS();
    cy.contains('#navbarCollapse ul.navbar-nav li.nav-item a.nav-link', 'example@example.com').click();
    cy.contains('#navbarCollapse ul.navbar-nav li.nav-item a.nav-link', 'example@example.com').parent().find('.dropdown-menu').contains('a.dropdown-item', 'My profile').click();
    cy.get('#sys-modal-user-settings').should('exist');
    cy.get('#sys-modal-user-settings form input[name=first_name]').should('have.value', 'FirstNew');
    cy.get('#sys-modal-user-settings form input[name=last_name]').should('have.value', 'LastNew');
    cy.get('#sys-modal-user-settings form button[type=button].btn').click();
    cy.logoutCMS();
  });
});
