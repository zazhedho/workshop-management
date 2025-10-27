import React, { useState } from 'react'
import { Outlet } from 'react-router-dom'
import { Container, Row, Col } from 'react-bootstrap'
import Sidebar from './Sidebar'
import Header from './Header'

const Layout = () => {
  const [sidebarOpen, setSidebarOpen] = useState(false)

  const toggleSidebar = () => {
    setSidebarOpen(!sidebarOpen)
  }

  const currentYear = new Date().getFullYear();
  const copyrightYear = currentYear > 2018 ? `2018-${currentYear}` : '2018';

  return (
    <Container fluid className="p-0">
      <Row className="g-0" style={{ minHeight: '100vh' }}>
        <Col md={3} lg={2} className={`sidebar ${sidebarOpen ? 'show' : ''}`}>
          <Sidebar />
        </Col>
        <Col md={9} lg={10} className="main-content d-flex flex-column">
          <Header toggleSidebar={toggleSidebar} />
          <Container fluid className="p-4 flex-grow-1">
            <Outlet />
          </Container>
          <footer className="footer mt-auto py-3 bg-light">
            <div className="container text-center">
              <span className="text-muted">© {copyrightYear} ZZ Family | <a href="https://www.linkedin.com/in/zaidus-zhuhur/" target="_blank" rel="noopener noreferrer">Help Center</a></span>
            </div>
          </footer>
        </Col>
      </Row>
    </Container>
  )
}

export default Layout