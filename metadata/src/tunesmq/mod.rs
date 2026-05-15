use anyhow::Result;

pub trait MqService {
    fn start(&mut self) -> Result<()>;
    fn stop(&mut self) -> Result<()>;
}

pub struct MqManager<'a> {
    services: Vec<Box<dyn MqService + 'a>>,
}

impl<'a> MqManager<'a> {
    pub fn new() -> Self {
        Self { services: vec![] }
    }

    pub fn start(&mut self) -> Result<()> {
        for service in &mut self.services {
            service.start()?;
        }

        Ok(())
    }

    pub fn stop(&mut self) -> Result<()> {
        for service in &mut self.services {
            service.stop()?;
        }

        Ok(())
    }

    pub fn register(&mut self, service: impl MqService + 'a) -> &mut Self {
        self.services.push(Box::new(service));
        self
    }
}
