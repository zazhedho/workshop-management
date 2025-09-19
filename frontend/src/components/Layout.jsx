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

  return (
    <Container fluid className="p-0">
      <Row className="g-0">
        <Col md={3} lg={2} className={`sidebar ${sidebarOpen ? 'show' : ''}`}>
          <Sidebar />
        </Col>
        <Col md={9} lg={10} className="main-content">
          <Header toggleSidebar={toggleSidebar} />
          <Container fluid className="p-4">
            <Outlet />
          </Container>
        </Col>
      </Row>
    </Container>
  )
}

export default Layout