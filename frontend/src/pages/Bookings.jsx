import React, { useState, useEffect } from 'react'
import { Card, Table, Button, Modal, Form, Row, Col, Pagination, Alert, Badge } from 'react-bootstrap'
import api from '../services/api'
import { useAuth } from '../contexts/AuthContext'

const Bookings = () => {
  const { user } = useAuth()
  const [bookings, setBookings] = useState([])
  const [vehicles, setVehicles] = useState([])
  const [services, setServices] = useState([])
  const [loading, setLoading] = useState(true)
  const [showModal, setShowModal] = useState(false)
  const [showStatusModal, setShowStatusModal] = useState(false)
  const [selectedBooking, setSelectedBooking] = useState(null)
  const [currentPage, setCurrentPage] = useState(1)
  const [totalPages, setTotalPages] = useState(1)
  const [search, setSearch] = useState('')
  const [filters, setFilters] = useState({
    status: ''
  })
  const [error, setError] = useState('')
  const [success, setSuccess] = useState('')

  const [formData, setFormData] = useState({
    vehicle_id: '',
    booking_date: '',
    notes: '',
    service_ids: []
  })

  const [statusData, setStatusData] = useState({
    status: ''
  })

  useEffect(() => {
    fetchBookings()
    fetchVehicles()
    fetchServices()
  }, [currentPage, search, filters])

  const fetchBookings = async () => {
    try {
      setLoading(true)
      const params = new URLSearchParams({
        page: currentPage,
        limit: 10,
        ...(search && { search }),
        ...(filters.status && { 'filters[status]': filters.status })
      })

      const response = await api.get(`/bookings?${params}`)
      setBookings(response.data.data || [])
      setTotalPages(response.data.total_pages || 1)
    } catch (error) {
      console.error('Failed to fetch bookings:', error)
      setError('Failed to fetch bookings')
    } finally {
      setLoading(false)
    }
  }

  const fetchVehicles = async () => {
    try {
      const response = await api.get('/vehicles?limit=100')
      setVehicles(response.data.data || [])
    } catch (error) {
      console.error('Failed to fetch vehicles:', error)
    }
  }

  const fetchServices = async () => {
    try {
      const response = await api.get('/services?limit=100')
      setServices(response.data.data || [])
    } catch (error) {
      console.error('Failed to fetch services:', error)
    }
  }

  const handleShowModal = () => {
    setFormData({
      vehicle_id: '',
      booking_date: '',
      notes: '',
      service_ids: []
    })
    setShowModal(true)
    setError('')
    setSuccess('')
  }

  const handleCloseModal = () => {
    setShowModal(false)
    setError('')
    setSuccess('')
  }

  const handleShowStatusModal = (booking) => {
    setSelectedBooking(booking)
    setStatusData({ status: booking.status })
    setShowStatusModal(true)
    setError('')
    setSuccess('')
  }

  const handleCloseStatusModal = () => {
    setShowStatusModal(false)
    setSelectedBooking(null)
    setError('')
    setSuccess('')
  }

  const handleChange = (e) => {
    const { name, value, type, checked } = e.target
    
    if (name === 'service_ids') {
      if (checked) {
        setFormData(prev => ({
          ...prev,
          service_ids: [...prev.service_ids, value]
        }))
      } else {
        setFormData(prev => ({
          ...prev,
          service_ids: prev.service_ids.filter(id => id !== value)
        }))
      }
    } else {
      setFormData({
        ...formData,
        [name]: value
      })
    }
  }

  const handleStatusChange = (e) => {
    setStatusData({
      status: e.target.value
    })
  }

  const handleSubmit = async (e) => {
    e.preventDefault()
    setError('')
    setSuccess('')

    try {
      const submitData = {
        ...formData,
        booking_date: new Date(formData.booking_date).toISOString()
      }

      await api.post('/booking', submitData)
      setSuccess('Booking created successfully')
      
      fetchBookings()
      setTimeout(() => {
        handleCloseModal()
      }, 1500)
    } catch (error) {
      setError(error.response?.data?.error || 'Operation failed')
    }
  }

  const handleStatusSubmit = async (e) => {
    e.preventDefault()
    setError('')
    setSuccess('')

    try {
      await api.put(`/booking/${selectedBooking.id}/status`, statusData)
      setSuccess('Booking status updated successfully')
      
      fetchBookings()
      setTimeout(() => {
        handleCloseStatusModal()
      }, 1500)
    } catch (error) {
      setError(error.response?.data?.error || 'Status update failed')
    }
  }

  const handleSearch = (e) => {
    setSearch(e.target.value)
    setCurrentPage(1)
  }

  const handleFilterChange = (key, value) => {
    setFilters(prev => ({ ...prev, [key]: value }))
    setCurrentPage(1)
  }

  const getStatusBadge = (status) => {
    const variants = {
      pending: 'warning',
      confirmed: 'info',
      'on progress': 'primary',
      completed: 'success',
      cancelled: 'danger'
    }
    return <Badge bg={variants[status] || 'secondary'}>{status}</Badge>
  }

  const formatDate = (dateString) => {
    return new Date(dateString).toLocaleDateString('en-US', {
      year: 'numeric',
      month: 'short',
      day: 'numeric',
      hour: '2-digit',
      minute: '2-digit'
    })
  }

  const canCreate = user?.role === 'customer' || user?.role === 'admin'
  const canUpdateStatus = user?.role === 'admin' || user?.role === 'cashier'

  return (
    <div>
      <div className="d-flex justify-content-between align-items-center mb-4">
        <h2>Bookings Management</h2>
        {canCreate && (
          <Button variant="primary" onClick={handleShowModal}>
            <i className="fas fa-plus me-2"></i>
            New Booking
          </Button>
        )}
      </div>

      {error && <Alert variant="danger">{error}</Alert>}
      {success && <Alert variant="success">{success}</Alert>}

      {/* Filters */}
      <Card className="mb-4">
        <Card.Body>
          <Row>
            <Col md={6}>
              <Form.Group>
                <Form.Label>Search Bookings</Form.Label>
                <Form.Control
                  type="text"
                  placeholder="Search by notes..."
                  value={search}
                  onChange={handleSearch}
                />
              </Form.Group>
            </Col>
            <Col md={3}>
              <Form.Group>
                <Form.Label>Filter by Status</Form.Label>
                <Form.Select
                  value={filters.status}
                  onChange={(e) => handleFilterChange('status', e.target.value)}
                >
                  <option value="">All Status</option>
                  <option value="pending">Pending</option>
                  <option value="confirmed">Confirmed</option>
                  <option value="on progress">On Progress</option>
                  <option value="completed">Completed</option>
                  <option value="cancelled">Cancelled</option>
                </Form.Select>
              </Form.Group>
            </Col>
          </Row>
        </Card.Body>
      </Card>

      {/* Bookings Table */}
      <Card>
        <Card.Header>
          <h5 className="mb-0">Bookings List</h5>
        </Card.Header>
        <Card.Body>
          {loading ? (
            <div className="text-center py-4">
              <div className="spinner-border" role="status">
                <span className="visually-hidden">Loading...</span>
              </div>
            </div>
          ) : bookings.length > 0 ? (
            <>
              <Table responsive hover>
                <thead>
                  <tr>
                    <th>Booking Date</th>
                    <th>Vehicle</th>
                    <th>Status</th>
                    <th>Notes</th>
                    <th>Created At</th>
                    {canUpdateStatus && <th>Actions</th>}
                  </tr>
                </thead>
                <tbody>
                  {bookings.map((booking) => (
                    <tr key={booking.id}>
                      <td>{formatDate(booking.booking_date)}</td>
                      <td>
                        <div className="d-flex align-items-center">
                          <i className="fas fa-calendar-check fa-2x text-primary me-3"></i>
                          <div>
                            <div className="fw-bold">{booking.vehicle_id}</div>
                          </div>
                        </div>
                      </td>
                      <td>{getStatusBadge(booking.status)}</td>
                      <td>{booking.notes || '-'}</td>
                      <td>{formatDate(booking.created_at)}</td>
                      {canUpdateStatus && (
                        <td>
                          <Button
                            variant="outline-primary"
                            size="sm"
                            onClick={() => handleShowStatusModal(booking)}
                          >
                            <i className="fas fa-edit me-1"></i>
                            Update Status
                          </Button>
                        </td>
                      )}
                    </tr>
                  ))}
                </tbody>
              </Table>

              {/* Pagination */}
              {totalPages > 1 && (
                <div className="d-flex justify-content-center mt-4">
                  <Pagination>
                    <Pagination.Prev
                      disabled={currentPage === 1}
                      onClick={() => setCurrentPage(currentPage - 1)}
                    />
                    {[...Array(totalPages)].map((_, index) => (
                      <Pagination.Item
                        key={index + 1}
                        active={index + 1 === currentPage}
                        onClick={() => setCurrentPage(index + 1)}
                      >
                        {index + 1}
                      </Pagination.Item>
                    ))}
                    <Pagination.Next
                      disabled={currentPage === totalPages}
                      onClick={() => setCurrentPage(currentPage + 1)}
                    />
                  </Pagination>
                </div>
              )}
            </>
          ) : (
            <p className="text-muted text-center py-4">No bookings found</p>
          )}
        </Card.Body>
      </Card>

      {/* Add Booking Modal */}
      <Modal show={showModal} onHide={handleCloseModal} size="lg">
        <Modal.Header closeButton>
          <Modal.Title>Create New Booking</Modal.Title>
        </Modal.Header>
        <Form onSubmit={handleSubmit}>
          <Modal.Body>
            {error && <Alert variant="danger">{error}</Alert>}
            {success && <Alert variant="success">{success}</Alert>}
            
            <Row>
              <Col md={6}>
                <Form.Group className="mb-3">
                  <Form.Label>Vehicle *</Form.Label>
                  <Form.Select
                    name="vehicle_id"
                    value={formData.vehicle_id}
                    onChange={handleChange}
                    required
                  >
                    <option value="">Select Vehicle</option>
                    {vehicles.map((vehicle) => (
                      <option key={vehicle.id} value={vehicle.id}>
                        {vehicle.license_plate} - {vehicle.brand} {vehicle.model}
                      </option>
                    ))}
                  </Form.Select>
                </Form.Group>
              </Col>
              <Col md={6}>
                <Form.Group className="mb-3">
                  <Form.Label>Booking Date *</Form.Label>
                  <Form.Control
                    type="datetime-local"
                    name="booking_date"
                    value={formData.booking_date}
                    onChange={handleChange}
                    required
                  />
                </Form.Group>
              </Col>
            </Row>

            <Form.Group className="mb-3">
              <Form.Label>Services *</Form.Label>
              <div className="border rounded p-3" style={{ maxHeight: '200px', overflowY: 'auto' }}>
                {services.map((service) => (
                  <Form.Check
                    key={service.id}
                    type="checkbox"
                    id={`service-${service.id}`}
                    name="service_ids"
                    value={service.id}
                    label={`${service.name} - $${service.price}`}
                    onChange={handleChange}
                    className="mb-2"
                  />
                ))}
              </div>
            </Form.Group>

            <Form.Group className="mb-3">
              <Form.Label>Notes</Form.Label>
              <Form.Control
                as="textarea"
                rows={3}
                name="notes"
                value={formData.notes}
                onChange={handleChange}
                placeholder="Additional notes..."
              />
            </Form.Group>
          </Modal.Body>
          <Modal.Footer>
            <Button variant="secondary" onClick={handleCloseModal}>
              Cancel
            </Button>
            <Button variant="primary" type="submit">
              Create Booking
            </Button>
          </Modal.Footer>
        </Form>
      </Modal>

      {/* Update Status Modal */}
      <Modal show={showStatusModal} onHide={handleCloseStatusModal}>
        <Modal.Header closeButton>
          <Modal.Title>Update Booking Status</Modal.Title>
        </Modal.Header>
        <Form onSubmit={handleStatusSubmit}>
          <Modal.Body>
            {error && <Alert variant="danger">{error}</Alert>}
            {success && <Alert variant="success">{success}</Alert>}
            
            <Form.Group className="mb-3">
              <Form.Label>Status *</Form.Label>
              <Form.Select
                name="status"
                value={statusData.status}
                onChange={handleStatusChange}
                required
              >
                <option value="pending">Pending</option>
                <option value="confirmed">Confirmed</option>
                <option value="on progress">On Progress</option>
                <option value="completed">Completed</option>
                <option value="cancelled">Cancelled</option>
              </Form.Select>
            </Form.Group>
          </Modal.Body>
          <Modal.Footer>
            <Button variant="secondary" onClick={handleCloseStatusModal}>
              Cancel
            </Button>
            <Button variant="primary" type="submit">
              Update Status
            </Button>
          </Modal.Footer>
        </Form>
      </Modal>
    </div>
  )
}

export default Bookings