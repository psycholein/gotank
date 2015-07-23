window.Motion = class Motion extends window.BasicModule
  config: =>
    @position = "right"

  router: (event) =>
    switch event.Task
      when "motion"
        $module = $(@selector).filter("[data-name='#{event.Name}']")
        $module.find('.value').text(event.Data.value)
