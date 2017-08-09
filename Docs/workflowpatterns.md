# Workflow Patterns Support

This document describes how to implement the common workflow patterns as by the [Workflow Pattern Initiative](http://www.workflowpatterns.com/).
  
  
### Control Patterns
The control-flow perspective captures aspects related to control-flow dependencies between various tasks (e.g. parallelism, choice, synchronization etc).

#### 1. Sequence
- Description: Task B is enabled after the completion of preceding task A.
- Status: **Supported**
- Example: 
```json
{
  "taskA" : {
    "type" : "function",
    "name" : "funcA"
  },
  "taskB" : {
    "type" : "function",
    "name" : "funcB",
    "dependencies" : [
      "taskA"
    ]
  }
}
```

#### 2. Parallel Split
- Description: divergence of a branch A into two or more parallel branches executing concurrently (in this case B and C)
- Status: **Supported**
- Example:
```json
{
  "taskA" : {
    "type" : "function",
    "name" : "funcA"
  },
  "taskB" : {
    "type" : "function",
    "name" : "funcB",
    "dependencies" : [
      "taskA"
    ]
  },
  "taskC" : {
    "type" : "function",
    "name" : "funcC",
    "dependencies" : [
      "taskA"
    ]
  }
}
```


#### 3. Synchronization
- Description: convergence of two or more branches (in this case A and B) into a single branch C
- Status: **Supported (though task C will receive only 1 of the outputs)**
- Example:
```json
{
  "taskA" : {
    "type" : "function",
    "name" : "funcA"
  },
  "taskB" : {
    "type" : "function",
    "name" : "funcB",
    "dependencies" : []
  },
  "taskC" : {
    "type" : "function",
    "name" : "funcC",
    "dependencies" : [
      "taskA", 
      "TaskB"
    ]
  }
}
```

#### 4. Exclusive Choice (XOR / Switch) 
- Description: divergence of branch A into one of the two or more branches based on the mechanism
- Status: **Not Supported**

#### 6. Multi-choice
- Description: divergence of branch A into one or more of the two or more branches based on the mechanism
- Status: **Not Supported**

#### 7. Structured Synchronizing Merge
- Description: A synchronization in which branches that are not activated (e.g. with a multi-choice) are not waited for.
- Status: **Not Supported**
- Note: This can be supported implicitly by propagating a SKIP token through non-choice branches. 

#### 8. Multi-Merge
- Description: A task with multiple inputs that is activated on each input activation. 
- Status: **Not Supported**
- Note: question of whether it will ever be supported. It relies on multiple link activations

#### 9. Structured Discriminator (Race) and 28. Blocking Discriminator
- Description: A task with multiple inputs that is activated only for the first link activation that occurs. Subsequent activations are ignored.
- Status: **Not Supported**

#### 10. Arbitrary Cycles
- Description: Ability to represent cycles within a workflow 
- Status: **Partially Supported**
- Note: This will probably never be supported explicitly. Instead recursion can be used. However, some shortcuts might be useful.

#### 11. Implicit Termination
- Description: A workflow invocation should terminate when there are no remaining tasks to be done now or at any point in the future. (no deadlock)
- Status: **Supported**
- Note: Deadlocks are not possible currently. The workflow invocation completes when no tasks are remaining.

#### 12. Multiple Instances without Synchronization (duplicate task invocations)
- Description: A task can instantiated multiple times and can run concurrently, in their own context.
- Status: **(Probably) Never Supported(?)**
- Note: Use separate task in branches or use recursion.

#### 13. Multiple Instances with a priori Design-Time Knowledge
- Description: Multiple, independent, concurrent instances of a task can be created, where the number of instances is known at design time. Before continuing the tasks need to be synchronized.
- Status: **(Probably) Never Supported(?)**
- Note: Use duplicate tasks (loop unrolling) or use recursion.

#### 14. Multiple Instances with a priori Run-Time Knowledge
- Description: Multiple, independent, concurrent instances of a task can be created, where the number of instances is known at workflow invocation run-time, because of state data, resource availability, communication, etc. 
Before continuing the tasks need to be synchronized.
- Status: **(Probably) Never Supported(?)**
- Note: Use recursion.

#### 15. Multiple Instances without a priori Run-Time Knowledge
- Description: Multiple, independent, concurrent instances of a task can be created, where the number of instances is not known until the last instance has completed, because of state data, resource availability, communication, etc. 
Before continuing the tasks need to be synchronized.
- Status: **(Probably) Never Supported(?)**
- Note: use recursion.

#### 16. Deferred Choice
- Description: The first task of several branches runs, based on the result of those tasks branch(es) are chosen.
- Status: **Not Supported**
- Note: Seems overlapping with 9. 

#### 17. Interleaved Parallel Routing
- Description: Tasks have relaxed/partial/no ordering, but cannot be done concurrently.
- Status: **Not Supported**
- Note: Does not seem to add much, for the perceived complexity.

#### 18. Milestone
- Description: A task is only activated when the workflow invocation is in a specific state (commonly in a parallel branch). If the workflow invocation has already progressed beyond the specific task, the task is not invoked. (deadline passed)
- Status: **Not Supported**
- Note: Temporal link activation.

#### 19. Cancel Task
- Description: An link activation is canceled or (if supported) the task is halted/aborted.
- Status: **Not Supported**
- Note: High prio

#### 20. Cancel Case (Cancel Worfklow Invocation)
- Description: Cancel the entire workflow Invocation, canceling all tasks.
- Status: **Not Supported**
- Note: High prio

#### 21. Structured Loop
- Description: loop with a (post- or pre-)condition
- Status: **Not Supported**
- Note: This will probably never be supported explicitly. Instead recursion can be used. However, some shortcuts might be useful.

#### 22. Recursion
- Description: ability to calling itself or any other parent workflows
- Status: **Supported**

#### 23. Transient Trigger
- Description: The ability for a task to be triggered by a signal from another part of the invocation or from the external environment. 
These triggers are transient in nature and are lost if not acted on immediately by the receiving task. 
A trigger can only be utilized if there is a task instance waiting for it at the time it is received.
- Status: **Not Supported**
- Note: Not sure if this should be supported explicitly. Could be implemented using a watch do...while/recursive, delayed task.

#### 24. Persistent Trigger
- Description: The ability for a task to be triggered by a signal from another part of the process or from the external environment. 
These triggers are persistent in form and are retained by the process until they can be acted on by the receiving task.
- Status: **Not Supported**
- Note: Not sure if this should be supported explicitly. Could be implemented using a watch do...while/recursive, delayed, starting task.

#### 25. Cancel Region
- Description: The ability to disable, cancel a set of tasks.
- Status: **Not Supported**
- Note: Could either be supported by supporting closures (unnamed workflow in workflow), or just solved using sub-workflow.

#### 26. Cancel Multiple Instance Task, and 27. Complete Multiple Instance Task 
- Description: Cancel all duplicate task invocations
- Status: **Not Supported**
- Note: As tasks are probably only allowed to be invoked once, this will not be needed.
 
#### 29. Cancelling Discriminator 
- Description: Similar to discriminator, only here once the discrimator is enabled branches other than the enabling one are canceled.
- Status: **Not Supported**

#### 30. Structured Partial Join and 31. Blocking Partial Join and 32. Canceling Partial Join
- Description: Similar to a discriminator, only here the subsequent link is enabled only when n of the incoming links is enabled.
- Status: **Not Supported**

#### 33. Generalized AND-join
- Description: Similar to synchronization, however once a link is enabled this memory persists for the AND condition.
- Status: **Supported**
- Note: There is no way that a link can be disabled after it has been enabled currently.

#### 34. Static Partial Join for Multiple Instances and 35. Canceling Partial Join for Multiple Instances
- Description: Continue if N < M duplicate task invocations have completed. Subsequent ones are executed but ignored.
- Status: **Not Supported**

#### 36. Dynamic Partial Join for Multiple Instances
- Description: Similar to 34, only in this case the number of total tasks is not known, so whether to proceed is determined by a condition.
- Status: **Not Supported**

#### 37. Local Synchronizing Merge
- Description: Similar to a merge, but the decision on which branches to wait is constantly evaluated.
- Status: **Not Supported** 

#### 38. General Synchronizing Merge
- Description: Similar to a merge, but here the next link is enabled either when all possible incoming links are enabled or when none can and are enabled.
- Status: **Not Supported**

#### 39. Critical Section
- Description: Given a workflow with two critical sections, when execution in one critical section starts, the other cannot start after the first one completes.
- Status: **Not Supported**
- Note: can be implemented implicitly using closures/sub-workflows

#### 40. Interleaved Routing
- Description: Tasks need to be executed in any order, as long is it is not occurring concurrently.
- Status: **Not Supported**
- Note: Might require mutex functionality.

### Data Patterns
The data perspective deals with the passing of information , scoping of variables, etc

#### 1. Task Data
- Description: Static data that is only available to a specific task. Defined as a task parameter or inside the task definition itself
- Status: **Not Supported**

#### 2. Block Data
- Description: Block tasks (i.e. tasks which can be described in terms of a corresponding subprocess) are able to define 
data elements which are accessible by each of the components of the corresponding sub-invocation.
- Status: **Not Supported**

#### 3. Scope Data
- Description: Data elements can be defined which are accessible by a subset of the tasks in a case.
- Status: **Not Supported**
- Note: Could incorporate task data

#### 4. Multiple Instance Data
- Description: Tasks which are able to execute multiple times within a single case can define data elements which are specific to an individual execution instance.
- Status: **Not/never Supported**
- Note: depends on whether to support duplicate task invocations.

#### 5. Case Data
- Description: Data associated with a specific workflow invocation.
- Status **Not Supported**

#### 6. Folder Data
- Description: Data stored in folders that can be selectively fetched by workflow invocations
- Status: **Never Supported**
- Note: Seems out of scope.

#### 7. Workflow Data
- Description: Data elements are supported which are accessible to all components in each and every case of the process and are within the context of the process itself.
- Status: **Not supported**

#### 8. Environment Data
- Description: Data elements which exist in the external operating environment are able to be accessed by components of processes during execution.
- Status: **Supported**
- Note: Simply uses env.

#### 9. Task to Task
#### 10. Block to Sub-workflow
#### 11. Sub-workflow to block
- Description: The ability to communicate data elements between one task instance and another within the same case.
- Status: **Partially Supported** (only 1-to-1)
- Note: approaches possible: integrated communication, distinct data and control lines and shared data store. Pass by value, pass by reference (location) possible.

#### 12. To Multiple Task
- Description: The ability to pass data to multiple task invocations
- Status: **Not Supported**

#### 13. From Multiple Task
- Description: The ability to collect/aggregate the data of the tasks into a single message, and sending it to the receiver.

#### 14. Case to Case (Workflow Invocation to Workflow Invocation)
- Description: passing data to a concurrently running workflow invocation.
- Status: **Not Supported**
- Note: not clear why every needed.

#### 15. Task to Environment - Push
#### 16. Environment to Task - Pull/Push
- Description: The ability of a task to interact with data elements to and from a resource or service in the operating environment.
- Status: **Mostly never supported or supported using a function**
- Note: point 15 until 26. Can interact  


#### 27. Data Transfer by Value - Incoming
- Description: The ability of a process component to receive incoming data elements by value avoiding the need to have shared names or common address space with the component(s) from which it receives them.
- Status: **Supported**

#### 28. Data Transfer by Value - Outgoing
- Description: The ability of a process component to pass data elements to subsequent components as values avoiding the need to have shared names or common address space with the component(s) to which it is passing them.
- Status: **Supported**

#### 29. Data Transfer - Copy In/Copy Out
- Description: The ability of a process component to copy the values of a set of data elements from an external source (either within or outside the process environment) into its address space at the commencement of execution and to copy their final values back at completion.
- Status: **Indirectly Supported**
- Note: Use functions to facilitate this.

#### 30. Data Transfer by Reference - Unlocked
- Description: The ability to communicate data elements between process components by utilizing a reference to the location of the data element in some mutually accessible location. No concurrency restrictions apply to the shared data element.
- Status: **Indirectly Supported**
- Note: Use functions to facilitate this.

#### 31. Data Transfer by Reference - With Lock
- Description: The ability to communicate data elements between process components by passing a reference to the location of the data element in some mutually accessible location.
- Status: **Indirectly Supported**       
- Note: Use functions to facilitate this.

#### 32. Data Transformation - Input
- Description: The ability to apply a transformation function to a data element prior to it being passed to a process component. 
The transformation function has access to the same data elements as the receiving process component.
- Status: **Not Supported**

#### 33. Data Transformation - Output
- Description: The ability to apply a transformation function to a data element immediately prior to it being passed out of a process component. 
The transformation function has access to the same data elements as the process component that initiates it.
- Status: **Not Supported**

#### 34. Task Precondition - Data Existence
- Description: Data-based preconditions can be specified for tasks based on the presence of data elements at the time of execution. 
The preconditions can utilize any data elements available to the task with which they are associated. 
A task can only proceed if the associated precondition evaluates positively.       
- Status: **Supported**

#### 35. Task Precondition - Data Value
- Description: Data-based preconditions can be specified for tasks based on the value of specific parameters at the time of execution. 
The preconditions can utilize any data elements available to the task with which they are associated. 
A task can only proceed if the associated precondition evaluates positively.
- Status: **Not Supported**


#### 36. Task Postcondition - Data Existence
#### 37. Task Postcondition - Data Value
- Description: Data-based postconditions can be specified for tasks based on the existence of specific parameters at the time of task completion. 
The postconditions can utilize any data elements available to the task with which they are associated. 
A task can only proceed if the associated postcondition evaluates positively.
- Status: **Not Supported**
- Note: this would require either re-running the task, or delay execution until condition is satisfied.

#### 38. Event-Based Task Trigger
- Description: The ability for an external event to initiate a task and to pass data elements to it.
- Status: **Not Supported**

#### 39. Data-Based Task Trigger
- Description: Data-based task triggers provide the ability to trigger a specific task when an expression based on data elements in the process instance evaluates to true. 
- Status: **Not Supported**

#### 40. Data-based routing
- Description: Data-based routing provides the ability to alter the control-flow within a case based on the evaluation of data-based expressions. 
A data-based routing expression is associated with each outgoing arc of an OR-split or XOR-split.
- Status: **Not Supported**
               
### Resource Patterns
The resource perspective aims to capture the various ways in which resources are represented and utilized in workflows.



### Exception Handling Patterns
The patterns for the exception handling perspective deal with the various causes of exceptions and the various actions that need to be taken as a result of exceptions occurring.

### Notes
- Process = workflow invocation
- Work item = task
